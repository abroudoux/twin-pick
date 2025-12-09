package scrapper

import (
	"time"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/charmbracelet/log"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// BrowserScrapper implements domain.SuggestionsProvider using a headless browser
// to handle JavaScript-rendered content on Letterboxd's popular films page
type BrowserScrapper struct{}

func NewBrowserScrapper() *BrowserScrapper {
	return &BrowserScrapper{}
}

// GetSuggestions implements domain.SuggestionsProvider
func (b *BrowserScrapper) GetSuggestions(params *domain.ScrapperFilters) ([]*domain.Film, error) {
	url := buildPopularFilmsURL(params)
	return b.scrapeFilmsWithBrowser(url)
}

func (b *BrowserScrapper) scrapeFilmsWithBrowser(url string) ([]*domain.Film, error) {
	log.Infof("Launching browser to scrape: %s", url)

	l := launcher.New().
		Headless(true).
		NoSandbox(true).
		Set("disable-gpu").
		Set("disable-software-rasterizer")
	defer l.Cleanup()

	controlURL, err := l.Launch()
	if err != nil {
		return nil, err
	}

	browser := rod.New().ControlURL(controlURL)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
	defer browser.MustClose()

	page := browser.MustPage("").Timeout(60 * time.Second)

	// Set a realistic user agent before navigation
	page.MustSetUserAgent(&proto.NetworkSetUserAgentOverride{
		UserAgent:      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		AcceptLanguage: "en-US,en;q=0.9",
		Platform:       "Linux",
	})

	// Navigate to the page
	if err := page.Navigate(url); err != nil {
		return nil, err
	}

	// Wait for network to be idle (no requests for 500ms)
	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	// Wait for idle state
	page.MustWaitIdle()

	// Wait for the poster grid to be rendered by JavaScript
	// Based on letterboxdpy selectors: li.posteritem, li.griditem, div.react-component
	selectors := []string{
		"li.posteritem",
		"li.griditem",
		"li.poster-container",
	}

	var foundSelector string
	for _, selector := range selectors {
		_, err = page.Timeout(10 * time.Second).Element(selector)
		if err == nil {
			foundSelector = selector
			log.Infof("Found elements with selector: %s", selector)
			break
		}
	}

	if foundSelector == "" {
		log.Warnf("No poster elements found with any selector")
		return nil, nil
	}

	// Give more time for all elements to render
	time.Sleep(5 * time.Second)

	elements, err := page.Elements(foundSelector)
	if err != nil {
		return nil, err
	}

	log.Infof("Found %d elements to process", len(elements))

	var films []*domain.Film

	for _, el := range elements {
		// Method 1: div.react-component with data-film-slug (from letterboxdpy)
		reactDiv, err := el.Element("div.react-component")
		if err == nil {
			// Try data-film-slug first, then data-item-slug
			slug, _ := reactDiv.Attribute("data-film-slug")
			if slug == nil || *slug == "" {
				slug, _ = reactDiv.Attribute("data-item-slug")
			}

			if slug != nil && *slug != "" {
				// Try data-item-name first
				title, _ := reactDiv.Attribute("data-item-name")
				if title == nil || *title == "" {
					// Fallback to img alt
					img, imgErr := reactDiv.Element("img")
					if imgErr == nil {
						alt, _ := img.Attribute("alt")
						if alt != nil {
							title = alt
						}
					}
				}

				if title != nil && *title != "" {
					detailsEndpoint := "/film/" + *slug + "/details/"
					films = append(films, domain.NewFilm(*title, detailsEndpoint))
					continue
				}
			}
		}

		// Method 2: div.film-poster as fallback
		div, err := el.Element("div.film-poster")
		if err == nil {
			slug, _ := div.Attribute("data-film-slug")
			if slug != nil && *slug != "" {
				img, imgErr := div.Element("img")
				var title string
				if imgErr == nil {
					alt, _ := img.Attribute("alt")
					if alt != nil {
						title = *alt
					}
				}
				if title != "" {
					detailsEndpoint := "/film/" + *slug + "/details/"
					films = append(films, domain.NewFilm(title, detailsEndpoint))
				}
			}
		}
	}

	log.Infof("Found %d popular films", len(films))
	return films, nil
}
