package application

import "github.com/abroudoux/twinpick/internal/domain"

type fakeSuggestionsProvider struct {
	results []domain.Film
	err     error
}

func (f *fakeSuggestionsProvider) GetSuggestions(_ *domain.ScrapperParams) ([]domain.Film, error) {
	if err, ok := f.err, f.err != nil; ok {
		return nil, err
	}
	return f.results, nil
}
