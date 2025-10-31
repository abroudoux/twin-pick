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

func (s *PickService) Pick(pp *domain.ProgramParams) ([]domain.Film, error) {
	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		watchlists = make(map[string]*domain.Watchlist)
		firstError error
	)

	for _, user := range pp.Usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()

			wl, err := s.WatchlistProvider.GetWatchlist(username, pp.ScrapperParams)
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

	commonFilms, err := domain.CompareWatchlists(watchlists)
	if err != nil {
		return nil, err
	}

	if pp.Limit > 0 && len(commonFilms) > pp.Limit {
		commonFilms = commonFilms[:pp.Limit]
	}

	return commonFilms, nil
}
