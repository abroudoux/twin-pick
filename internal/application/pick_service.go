package application

import (
	"sync"

	"github.com/abroudoux/twinpick/internal/domain"
)

type PickService struct {
	WatchlistProvider domain.WatchlistProvider
}

func NewPickService(wp domain.WatchlistProvider) *PickService {
	return &PickService{WatchlistProvider: wp}
}

func (s *PickService) Pick(usernames []string, sp *domain.ScrapperParams, limit int) ([]domain.Film, error) {
	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		watchlists = make(map[string]*domain.Watchlist)
		firstError error
	)

	for _, user := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()

			wl, err := s.WatchlistProvider.GetWatchlist(username, sp)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = err
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			watchlists[username] = wl
			mu.Unlock()
		}(user)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	commonFilms, err := domain.GetCommonFilms(watchlists)
	if err != nil {
		return nil, err
	}

	if limit > 0 && len(commonFilms) > limit {
		commonFilms = commonFilms[:limit]
	}

	return commonFilms, nil
}
