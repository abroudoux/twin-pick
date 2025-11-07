package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/abroudoux/twinpick/internal/domain"
)

type filmDetailsResponse struct {
	Result      bool `json:"result"`
	ReleaseYear int  `json:"releaseYear"`
	RunTime     int  `json:"runTime"`
	Directors   []struct {
		Name string `json:"name"`
	} `json:"directors"`
}

func GetFilmsDetails(films []*domain.Film) ([]*domain.Film, error) {
	if films == nil {
		return nil, nil
	}

	baseURL := "https://letterboxd.com"
	var wg sync.WaitGroup

	for _, film := range films {
		if film.DetailsEndpoint == "" {
			continue
		}

		wg.Add(1)
		go func(film *domain.Film) {
			defer wg.Done()
			fetchFilmDetails(film, baseURL)
		}(film)
	}

	wg.Wait()
	return films, nil
}

func fetchFilmDetails(film *domain.Film, baseURL string) {
	url := fmt.Sprintf("%s%s", baseURL, film.DetailsEndpoint)
	resp, err := http.Get(url)
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

	var directors []string
	for _, d := range details.Directors {
		directors = append(directors, d.Name)
	}
	film.Directors = directors
}
