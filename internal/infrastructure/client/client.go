package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"

	"github.com/abroudoux/twinpick/internal/domain"
)

const (
	baseURL            = "https://letterboxd.com"
	maxConcurrentFetch = 20
	requestTimeout     = 10 * time.Second
	userAgent          = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var httpClient = &http.Client{
	Timeout: requestTimeout,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
	},
}

type filmDetailsResponse struct {
	Result      bool `json:"result"`
	ReleaseYear int  `json:"releaseYear"`
	RunTime     int  `json:"runTime"`
	Directors   []struct {
		Name string `json:"name"`
	} `json:"directors"`
}

func GetFilmsDetails(films []*domain.Film) ([]*domain.Film, error) {
	if len(films) == 0 {
		return films, nil
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentFetch)

	for _, film := range films {
		if film.DetailsEndpoint == "" {
			continue
		}

		wg.Add(1)
		go func(f *domain.Film) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fetchFilmDetails(f)
		}(film)
	}

	wg.Wait()
	return films, nil
}

func fetchFilmDetails(film *domain.Film) {
	url := buildJSONEndpoint(film.DetailsEndpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Warnf("Failed to create request for %s: %v", film.Title, err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Warnf("Failed to fetch details for %s: %v", film.Title, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warnf("Non-200 status for %s: %d", film.Title, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warnf("Failed to read response for %s: %v", film.Title, err)
		return
	}

	var details filmDetailsResponse
	if err := json.Unmarshal(body, &details); err != nil {
		log.Warnf("Failed to unmarshal response for %s: %v", film.Title, err)
		return
	}

	film.Duration = details.RunTime
	film.Year = details.ReleaseYear

	directors := make([]string, 0, len(details.Directors))
	for _, d := range details.Directors {
		directors = append(directors, d.Name)
	}
	film.Directors = directors
}

// buildJSONEndpoint converts any film endpoint to the JSON endpoint
// /film/{slug}/details/ -> /film/{slug}/json/
// /film/{slug}/json/ -> /film/{slug}/json/ (unchanged)
// /film/{slug}/ -> /film/{slug}/json/
func buildJSONEndpoint(endpoint string) string {
	// If already a JSON endpoint, just prepend baseURL
	if strings.HasSuffix(endpoint, "/json/") {
		return baseURL + endpoint
	}

	endpoint = strings.TrimSuffix(endpoint, "details/")
	endpoint = strings.TrimSuffix(endpoint, "/")
	return fmt.Sprintf("%s%s/json/", baseURL, endpoint)
}
