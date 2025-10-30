package application

import (
	"github.com/abroudoux/twinpick/internal/domain"
)

type MatchService struct {
	Provider domain.WatchlistProvider
}

func NewMatchService(provider domain.WatchlistProvider) *MatchService {
	return &MatchService{Provider: provider}
}

func (s *MatchService) MatchFilm(usernames []string, params *domain.ScrapperParams) (domain.Film, error) {
	commonService := NewCommonService(s.Provider)
	commonFilms, err := commonService.GetCommonFilms(usernames, params)
	if err != nil {
		return domain.Film{}, err
	}

	return domain.SelectRandomFilm(commonFilms)
}
