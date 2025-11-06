package application

import (
	"sync"

	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/client"
)

type PickServiceInterface interface {
	Pick(pp *domain.PickParams) ([]*domain.Film, error)
}

type PickService struct {
	WatchlistProvider domain.WatchlistProvider
}

func NewPickService(wp domain.WatchlistProvider) *PickService {
	return &PickService{WatchlistProvider: wp}
}

func (s *PickService) Pick(pp *domain.PickParams) ([]*domain.Film, error) {
	var filmsDetailsFetched bool

	watchlists, err := s.collectWatchlists(pp.Usernames, pp.ScrapperParams)
	if err != nil {
		return nil, err
	}

	films, err := domain.CompareWatchlists(watchlists)
	if err != nil {
		return nil, err
	}

	if pp.Duration != domain.Long {
		films, err := client.FetchFilmsDetails(films)
		if err != nil {
			return nil, err
		}
		filmsDetailsFetched = true

		films = domain.FilterFilmsByDuration(films, pp.Duration)
	}

	if pp.Limit > 0 && len(films) > pp.Limit {
		films = films[:pp.Limit]
	}

	if !filmsDetailsFetched {
		films, err = client.FetchFilmsDetails(films)
		if err != nil {
			return nil, err
		}
	}

	return films, nil
}

func (s *PickService) collectWatchlists(usernames []string, params *domain.ScrapperParams) (map[string]*domain.Watchlist, error) {
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
			wl, err := s.WatchlistProvider.GetWatchlist(username, params)
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
