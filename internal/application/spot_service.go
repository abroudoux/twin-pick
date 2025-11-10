package application

import (
	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/client"
)

type SpotServiceInterface interface {
	Spot(spotParams *domain.SpotParams) ([]*domain.Film, error)
}

type SpotService struct {
	SuggestionsProvider domain.SuggestionsProvider
}

func NewSpotService(suggestionsProvider domain.SuggestionsProvider) *SpotService {
	return &SpotService{SuggestionsProvider: suggestionsProvider}
}

func (s *SpotService) Spot(spotParams *domain.SpotParams) ([]*domain.Film, error) {
	films, err := s.SuggestionsProvider.GetSuggestions(spotParams.Params.ScrapperFilters)
	if err != nil {
		return nil, err
	}

	films, err = client.GetFilmsDetails(films)
	if err != nil {
		return nil, err
	}

	return films, nil
}
