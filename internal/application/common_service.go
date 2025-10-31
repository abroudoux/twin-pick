package application

import (
	"sync"

	"github.com/charmbracelet/log"

	"github.com/abroudoux/twinpick/internal/domain"
)

type CommonService struct {
	Provider domain.WatchlistProvider
}

func NewCommonService(provider domain.WatchlistProvider) *CommonService {
	return &CommonService{Provider: provider}
}

func (s *CommonService) GetUsersWatchlists(usernames []string, params *domain.ScrapperParams) (map[string]*domain.Watchlist, error) {
	var (
		mu         sync.Mutex
		wg         sync.WaitGroup
		result     = make(map[string]*domain.Watchlist)
		firstError error
	)

	for _, user := range usernames {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()

			wl, err := s.Provider.GetWatchlist(username, params)
			if err != nil {
				mu.Lock()
				if firstError == nil {
					firstError = err
				}
				mu.Unlock()
				return
			}

			mu.Lock()
			result[username] = wl
			mu.Unlock()
		}(user)
	}

	wg.Wait()

	if firstError != nil {
		return nil, firstError
	}

	return result, nil
}

func (s *CommonService) GetCommonFilms(usernames []string, params *domain.ScrapperParams) ([]domain.Film, error) {
	watchlists, err := s.GetUsersWatchlists(usernames, params)
	if err != nil {
		return nil, err
	}

	commonFilms, err := domain.GetCommonFilms(watchlists)
	if err != nil {
		return nil, err
	}

	for _, f := range commonFilms {
		log.Infof("Common film: %s", f.Title)
	}

	return commonFilms, nil
}
