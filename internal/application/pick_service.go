package application

import (
	"sync"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/client"
)

type PickServiceInterface interface {
	Pick(pickParams *domain.PickParams) ([]*domain.Film, error)
}

type PickService struct {
	WatchlistProvider domain.WatchlistProvider
}

func NewPickService(watchlistProvider domain.WatchlistProvider) *PickService {
	return &PickService{WatchlistProvider: watchlistProvider}
}

func (s *PickService) Pick(pickParams *domain.PickParams) ([]*domain.Film, error) {
	watchlists, err := s.collectWatchlists(pickParams)
	if err != nil {
		return nil, err
	}

	films, err := domain.CompareWatchlists(watchlists, pickParams.Params.Filters.Strict)
	if err != nil {
		return nil, err
	}

	var detailsFetched bool

	if pickParams.Params.Filters.Duration != domain.Long {
		films, err := client.GetFilmsDetails(films)
		if err != nil {
			return nil, err
		}
		detailsFetched = true

		films = domain.FilterFilmsByDuration(films, pickParams.Params.Filters.Duration)
	}

	if pickParams.Params.Filters.Limit > 0 && len(films) > pickParams.Params.Filters.Limit {
		films = films[:pickParams.Params.Filters.Limit]
	}

	if detailsFetched {
		return films, nil
	}

	films, err = client.GetFilmsDetails(films)
	if err != nil {
		return nil, err
	}

	return films, nil
}

func (s *PickService) collectWatchlists(pickParams *domain.PickParams) (usersWatchlists map[string]*domain.Watchlist, err error) {
	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		watchlists = make(map[string]*domain.Watchlist)
		firstError error
	)

	for _, user := range pickParams.Usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			wl, err := s.WatchlistProvider.GetWatchlist(username, pickParams.Params.ScrapperFilters)
			mu.Lock()
			defer mu.Unlock()
			if err != nil && firstError == nil {
				firstError = err
				return
			}
			watchlists[username] = wl
		}(user)
	}
	wg.Wait()

	return watchlists, firstError
}
