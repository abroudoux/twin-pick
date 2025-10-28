package infrastructure

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/gocolly/colly/v2"
)

type LetterboxdScrapper struct{}

func NewLetterboxdScrapper() *LetterboxdScrapper {
	return &LetterboxdScrapper{}
}

func (s *LetterboxdScrapper) GetWatchlist(username string, params domain.ScrapperParams) (domain.Watchlist, error) {
	var films []domain.Film
	var totalPages int

	pageCollector := colly.NewCollector(colly.AllowedDomains("letterboxd.com"))
	pageCollector.OnHTML("div.paginate-pages ul", func(e *colly.HTMLElement) {
		e.ForEach("li.paginate-page a", func(_ int, el *colly.HTMLElement) {
			if n, err := strconv.Atoi(el.Text); err == nil && n > totalPages {
				totalPages = n
			}
		})
	})
	watchlistURL := buildWatchlistURL(username, params)
	_ = pageCollector.Visit(watchlistURL)
	pageCollector.Wait()

	filmCh := make(chan []domain.Film)
	var wg sync.WaitGroup
	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			c := colly.NewCollector(colly.AllowedDomains("letterboxd.com"))

			var pageFilms []domain.Film
			c.OnHTML("div.poster-grid li", func(e *colly.HTMLElement) {
				if name := e.ChildAttr("div.react-component", "data-item-full-display-name"); name != "" {
					pageFilms = append(pageFilms, domain.Film{Name: name})
				}
			})
			_ = c.Visit(fmt.Sprintf("%s/page/%d", strings.TrimRight(watchlistURL, "/"), page))
			c.Wait()
			filmCh <- pageFilms
		}(i)
	}

	go func() {
		wg.Wait()
		close(filmCh)
	}()

	for fs := range filmCh {
		films = append(films, fs...)
	}

	watchlist := domain.Watchlist{
		Films: films,
	}

	return watchlist, nil
}

func buildWatchlistURL(username string, params domain.ScrapperParams) string {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", username)
	if len(params.Genres) > 0 {
		url += "genre/" + strings.Join(params.Genres, "+") + "/"
	}
	if params.Platform != "" {
		url += fmt.Sprintf("on/%s/", params.Platform)
	}
	return url
}
