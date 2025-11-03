package scrapper

import (
	"strings"

	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *LetterboxdScrapper) GetSuggestions(params *domain.ScrapperParams) ([]domain.Film, error) {
	popularFilmsURL := buildPopularFilmsURL(params)
	return s.getPopularFilmsImpl(popularFilmsURL)
}

func (s *LetterboxdScrapper) getPopularFilmsImpl(popularFilmsURL string) ([]domain.Film, error) {
	var films []domain.Film

	collector := colly.NewCollector(colly.AllowedDomains("letterboxd.com"))
	collector.OnHTML("ul.poster-list li", func(e *colly.HTMLElement) {
		if title := e.ChildAttr("div.react-component", "data-item-full-display-name"); title != "" {
			films = append(films, domain.Film{Title: title})
		}
	})

	if err := collector.Visit(popularFilmsURL); err != nil {
		return nil, err
	}

	collector.Wait()
	return films, nil
}

func buildPopularFilmsURL(params *domain.ScrapperParams) string {
	url := "https://letterboxd.com/films/popular"

	if len(params.Genres) > 0 {
		url += "/genre/" + strings.Join(params.Genres, "+")
	}
	return url
}
