package application

import (
	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/abroudoux/twinpick/internal/infrastructure/client"
)

type SpotServiceInterface interface {
	Spot(sp *domain.SpotParams) ([]*domain.Film, error)
}

type SpotService struct {
	SuggestionsProvider domain.SuggestionsProvider
}

func NewSpotService(sp domain.SuggestionsProvider) *SpotService {
	return &SpotService{SuggestionsProvider: sp}
}

func (s *SpotService) Spot(sp *domain.SpotParams) ([]*domain.Film, error) {
	films, err := s.SuggestionsProvider.GetSuggestions(sp.ScrapperParams)
	if err != nil {
		return nil, err
	}

	filmsWithDetails, err := client.FetchFilmsDetails(films)
	if err != nil {
		return nil, err
	}

	return filmsWithDetails, nil
}
