package scrapper

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *LetterboxdScrapper) GetWatchlist(username string, params *domain.ScrapperParams) (*domain.Watchlist, error) {
	watchlist := domain.NewWatchlist(username)
	watchlistURL := buildWatchlistURL(username, params)

	totalPages, err := s.GetTotalWatchlistPages(watchlistURL)
	if err != nil {
		return nil, err
	}

	const maxConcurrent = 20
	pageCh := make(chan int, totalPages)
	filmCh := make(chan []domain.Film, totalPages)
	errCh := make(chan error, totalPages)

	for i := 1; i <= totalPages; i++ {
		pageCh <- i
	}
	close(pageCh)

	var wg sync.WaitGroup

	for range make([]struct{}, maxConcurrent) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for page := range pageCh {
				films, err := s.GetFilmsOnWatchlistPage(watchlistURL, page)
				if err != nil {
					errCh <- err
					return
				}
				filmCh <- films
			}
		}()
	}

	wg.Wait()
	close(filmCh)
	close(errCh)

	for fs := range filmCh {
		watchlist.Films = append(watchlist.Films, fs...)
	}

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return watchlist, nil
}

func (s *LetterboxdScrapper) getTotalWatchlistPagesImpl(watchlistURL string) (int, error) {
	totalPages := 1
	collector := s.NewCollector()
	collector.OnHTML("div.paginate-pages ul", func(e *colly.HTMLElement) {
		e.ForEach("li.paginate-page a", func(_ int, el *colly.HTMLElement) {
			if n, err := strconv.Atoi(el.Text); err == nil && n > totalPages {
				totalPages = n
			}
		})
	})

	if err := collector.Visit(watchlistURL); err != nil {
		return 0, err
	}
	collector.Wait()
	return totalPages, nil
}

func (s *LetterboxdScrapper) getFilmsOnWatchlistPageImpl(watchlistURL string, page int) ([]domain.Film, error) {
	var films []domain.Film

	collector := s.NewCollector()
	collector.OnHTML("div.poster-grid li", func(e *colly.HTMLElement) {
		if title := e.ChildAttr("div.react-component", "data-item-full-display-name"); title != "" {
			films = append(films, domain.Film{Title: title})
		}
	})

	pageURL := fmt.Sprintf("%s/page/%d", strings.TrimRight(watchlistURL, "/"), page)
	if err := collector.Visit(pageURL); err != nil {
		return nil, err
	}

	collector.Wait()
	return films, nil
}

func buildWatchlistURL(username string, params *domain.ScrapperParams) string {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist", username)

	if len(params.Genres) > 0 {
		url += "/genre/" + strings.Join(params.Genres, "+")
	}
	if params.Platform != "" {
		url += fmt.Sprintf("/on/%s", params.Platform)
	}
	return url
}
