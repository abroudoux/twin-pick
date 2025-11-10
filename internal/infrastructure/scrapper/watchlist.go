package scrapper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

const (
	maxConcurrentPages = 30
	totalTimeout       = 30 * time.Second
)

func (s *LetterboxdScrapper) GetWatchlist(username string, params *domain.ScrapperFilters) (*domain.Watchlist, error) {
	watchlist := domain.NewWatchlist(username)
	watchlistURL := buildWatchlistURL(username, params)

	totalPages, err := s.GetTotalWatchlistPages(watchlistURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get total pages: %w", err)
	}

	if totalPages == 0 {
		return watchlist, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel()

	estimatedFilms := totalPages * 28
	allFilms := make([]*domain.Film, 0, estimatedFilms)

	type pageResult struct {
		films []*domain.Film
		page  int
		err   error
	}

	resultCh := make(chan pageResult, totalPages)
	semaphore := make(chan struct{}, maxConcurrentPages)

	var wg sync.WaitGroup

	for page := 1; page <= totalPages; page++ {
		select {
		case <-ctx.Done():
			continue
		default:
		}

		wg.Add(1)
		go func(pageNum int) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				resultCh <- pageResult{err: ctx.Err(), page: pageNum}
				return
			}

			films, err := s.GetFilmsOnWatchlistPage(watchlistURL, pageNum)
			resultCh <- pageResult{
				films: films,
				page:  pageNum,
				err:   err,
			}
		}(page)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	errors := make([]error, 0)
	successCount := 0

	for result := range resultCh {
		if result.err != nil {
			errors = append(errors, fmt.Errorf("page %d: %w", result.page, result.err))
			continue
		}

		allFilms = append(allFilms, result.films...)
		successCount++
	}

	if len(errors) > 0 {
		if successCount == 0 {
			return nil, fmt.Errorf("all pages failed: %v", errors[0])
		}

		if len(errors) > totalPages/2 {
			return nil, fmt.Errorf("too many failures (%d/%d): %v",
				len(errors), totalPages, errors[0])
		}

		log.Warnf("Partial success: %d/%d pages retrieved, %d errors",
			successCount, totalPages, len(errors))
	}

	watchlist.Films = allFilms
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

func (s *LetterboxdScrapper) getFilmsOnWatchlistPageImpl(watchlistURL string, page int) ([]*domain.Film, error) {
	var films []*domain.Film

	collector := s.NewCollector()
	collector.OnHTML("div.poster-grid li", func(e *colly.HTMLElement) {
		title := e.ChildAttr("div.react-component", "data-item-full-display-name")
		detailsEndpoint := e.ChildAttr("div.react-component", "data-details-endpoint")

		if title == "" || detailsEndpoint == "" {
			return
		}

		films = append(films, domain.NewFilm(title, detailsEndpoint))
	})

	pageURL := fmt.Sprintf("%s/page/%d", strings.TrimRight(watchlistURL, "/"), page)
	if err := collector.Visit(pageURL); err != nil {
		return nil, err
	}

	collector.Wait()
	return films, nil
}
