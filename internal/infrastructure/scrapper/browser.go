package scrapper

import (
	"time"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/cache"
	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"
)

const cacheTTL = 24 * time.Hour

// BrowserScrapper implements domain.SuggestionsProvider
// Uses Letterboxd's AJAX endpoint to get popular films without needing a headless browser
type BrowserScrapper struct {
	cache *cache.Cache
}

func NewBrowserScrapper() *BrowserScrapper {
	return &BrowserScrapper{
		cache: cache.New(),
	}
}

// GetSuggestions implements domain.SuggestionsProvider
func (b *BrowserScrapper) GetSuggestions(params *domain.ScrapperFilters) ([]*domain.Film, error) {
	url := buildPopularFilmsAjaxURL(params)

	// Check cache first
	if cached, ok := b.cache.Get(url); ok {
		log.Infof("Cache hit for %s", url)
		return cached.([]*domain.Film), nil
	}

	// Scrape if not in cache
	films, err := b.scrapePopularFilms(url)
	if err != nil {
		return nil, err
	}

	// Store in cache with 24h TTL
	b.cache.Set(url, films, cacheTTL)
	log.Infof("Cached %d films for %s (TTL: %v)", len(films), url, cacheTTL)

	return films, nil
}

func (b *BrowserScrapper) scrapePopularFilms(url string) ([]*domain.Film, error) {
	log.Infof("Scraping popular films from: %s", url)

	var films []*domain.Film

	collector := colly.NewCollector()

	collector.OnHTML("li.posteritem", func(e *colly.HTMLElement) {
		// Same attributes as watchlist.go
		title := e.ChildAttr("div.react-component", "data-item-full-display-name")
		detailsEndpoint := e.ChildAttr("div.react-component", "data-details-endpoint")

		if title == "" || detailsEndpoint == "" {
			return
		}

		films = append(films, domain.NewFilm(title, detailsEndpoint))
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Errorf("Failed to scrape %s: %v", r.Request.URL, err)
	})

	if err := collector.Visit(url); err != nil {
		return nil, err
	}

	collector.Wait()

	log.Infof("Found %d popular films", len(films))
	return films, nil
}
