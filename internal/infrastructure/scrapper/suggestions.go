package scrapper

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/cache"
)

const popularFilmsTTL = 24 * time.Hour

var popularFilmsCache = cache.New()

func (s *LetterboxdScrapper) GetSuggestions(params *domain.ScrapperFilters) ([]*domain.Film, error) {
	popularFilmsURL := buildPopularFilmsURL(params)

	// Check cache first
	if cached, found := popularFilmsCache.Get(popularFilmsURL); found {
		films := cached.([]*domain.Film)
		log.Infof("Popular films cache hit: %d films", len(films))
		return films, nil
	}

	log.Infof("Popular films cache miss, scraping: %s", popularFilmsURL)
	films, err := s.getPopularFilmsImpl(popularFilmsURL)
	if err != nil {
		return nil, err
	}

	// Only cache if we found films (avoid caching empty results)
	if len(films) > 0 {
		popularFilmsCache.Set(popularFilmsURL, films, popularFilmsTTL)
		log.Infof("Cached %d popular films for 24h", len(films))
	} else {
		log.Warnf("No films found, not caching empty result")
	}

	return films, nil
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
