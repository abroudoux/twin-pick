package application

import (
	"fmt"

	"github.com/abroudoux/twinpick/internal/domain"
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

	for _, film := range films {
		fmt.Printf("Film suggested: %+v\n", film)
	}

	if sp.Limit > 0 && len(films) > sp.Limit {
		films = films[:sp.Limit]
	}

	return films, nil
}
