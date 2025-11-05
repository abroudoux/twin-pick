package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

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

func FetchFilmsDetails(films []*domain.Film) ([]*domain.Film, error) {
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

			url := fmt.Sprintf("%s%s", baseURL, film.DetailsEndpoint)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("[WARN] failed to fetch details for %s: %v", film.Title, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("[WARN] non-200 status for %s: %d", film.Title, resp.StatusCode)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("[WARN] failed to read response for %s: %v", film.Title, err)
				return
			}

			var details filmDetailsResponse
			if err := json.Unmarshal(body, &details); err != nil {
				log.Printf("[WARN] failed to unmarshal response for %s: %v", film.Title, err)
				return
			}

			film.Duration = details.RunTime
			film.Year = details.ReleaseYear

			var directors []string
			for _, d := range details.Directors {
				directors = append(directors, d.Name)
			}
			film.Directors = directors

		}(film)
	}

	wg.Wait()
	return films, nil
}
