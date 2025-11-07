package scrapper

import (
	"fmt"

	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *LetterboxdScrapper) GetFavouritesFilms(username string) ([]*domain.Film, error) {
	var favouritesFilms []*domain.Film

	for i := 4.0; i <= 5.0; i = i + 0.5 {
		films, _ := s.getFavouritesFilmsByNote(username, i)
		favouritesFilms = append(favouritesFilms, films...)
	}

	return favouritesFilms, nil
}

func (s *LetterboxdScrapper) getFavouritesFilmsByNote(username string, note float64) ([]*domain.Film, error) {
	if !isValidNote(note) {
		return nil, fmt.Errorf("note is not valid")
	}

	var films []*domain.Film

	collector := colly.NewCollector()
	collector.OnHTML("li.griditem div.react-component", func(e *colly.HTMLElement) {
		title := e.Attr("data-item-full-display-name")
		detailsEndpoint := e.Attr("data-details-endpoint")

		if title == "" || detailsEndpoint == "" {
			return
		}

		films = append(films, domain.NewFilm(title, detailsEndpoint))
	})

	favouritesFilmsURL := fmt.Sprintf("https://letterboxd.com/%s/films/rated/%f/by/date/", username, note)
	if err := collector.Visit(favouritesFilmsURL); err != nil {
		return nil, err
	}

	collector.Wait()
	return films, nil
}
