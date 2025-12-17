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

	// Apply limit before fetching details to optimize API calls
	limit := spotParams.Params.Filters.Limit
	if spotParams.Params.Filters.Duration != domain.Long {
		// If duration filter is applied, fetch more films then filter
		// to ensure we have enough results after filtering
		if limit > 0 {
			films = films[:min(limit*3, len(films))]
		}
	} else if limit > 0 && len(films) > limit {
		films = films[:limit]
	}

	films, err = client.GetFilmsDetails(films)
	if err != nil {
		return nil, err
	}

	// Apply duration filter after fetching details (need duration info)
	if spotParams.Params.Filters.Duration != domain.Long {
		films = domain.FilterFilmsByDuration(films, spotParams.Params.Filters.Duration)
	}

	// Apply final limit after duration filter
	if limit > 0 && len(films) > limit {
		films = films[:limit]
	}

	return films, nil
}
