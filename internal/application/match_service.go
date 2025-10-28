package application

import "github.com/abroudoux/twinpick/internal/domain"

type MatchService struct {
	Provider domain.WatchlistProvider
}

func NewMatchService(provider domain.WatchlistProvider) *MatchService {
	return &MatchService{Provider: provider}
}

func (s *MatchService) FindCommonFilm(usernames []string, params domain.ScrapperParams) (domain.Film, error) {
	watchlists := make(map[string]domain.Watchlist)
	for _, user := range usernames {
		wl, err := s.Provider.GetWatchlist(user, params)
		if err != nil {
			return domain.Film{}, err
		}
		watchlists[user] = wl
	}

	common, err := domain.GetCommonFilms(watchlists)
	if err != nil {
		return domain.Film{}, err
	}

	return domain.SelectRandomFilm(common)
}
