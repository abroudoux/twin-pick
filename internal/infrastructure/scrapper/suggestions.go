package scrapper

import (
	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *LetterboxdScrapper) GetSuggestions(params *domain.ScrapperParams) ([]*domain.Film, error) {
	popularFilmsURL := buildPopularFilmsURL(params)

	favouritesFilms, err := s.GetFavouritesFilms("abroudoux")
	if err != nil {
		return nil, err
	}
	log.Infof("%d", len(favouritesFilms))
	for _, f := range favouritesFilms {
		log.Infof("%s", f.Title)
	}

	return s.getPopularFilmsImpl(popularFilmsURL)
}

func (s *LetterboxdScrapper) getPopularFilmsImpl(popularFilmsURL string) ([]*domain.Film, error) {
	var films []*domain.Film

	collector := colly.NewCollector()
	collector.OnHTML("li.posteritem div.react-component", func(e *colly.HTMLElement) {
		title := e.Attr("data-item-full-display-name")
		detailsEndpoint := e.Attr("data-details-endpoint")

		if title == "" || detailsEndpoint == "" {
			return
		}

		films = append(films, domain.NewFilm(title, detailsEndpoint))
	})

	if err := collector.Visit(popularFilmsURL); err != nil {
		return nil, err
	}

	collector.Wait()
	return films, nil
}
