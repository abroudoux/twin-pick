package application

import (
	"errors"
	"reflect"
	"testing"

	"github.com/abroudoux/twinpick/internal/domain"
)

type fakeWatchlistProvider struct {
	results map[string]*domain.Watchlist
	errs    map[string]error
}

func (f *fakeWatchlistProvider) GetWatchlist(username string, _ *domain.ScrapperParams) (*domain.Watchlist, error) {
	if err, ok := f.errs[username]; ok {
		return nil, err
	}
	if wl, ok := f.results[username]; ok {
		return wl, nil
	}
	return nil, errors.New("user not found")
}

func makeWatchlist(username string, films ...string) *domain.Watchlist {
	wl := &domain.Watchlist{Username: username}
	for _, f := range films {
		wl.Films = append(wl.Films, domain.NewFilm(f, ""))
	}
	return wl
}

func TestCollectWatchlists_Success(t *testing.T) {
	provider := &fakeWatchlistProvider{
		results: map[string]*domain.Watchlist{
			"alice": makeWatchlist("alice", "Inception", "Matrix"),
			"bob":   makeWatchlist("bob", "Matrix", "Tenet"),
		},
	}
	service := NewPickService(provider)

	params := &domain.ScrapperParams{}
	watchlists, err := service.collectWatchlists([]string{"alice", "bob"}, params)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(watchlists) != 2 {
		t.Errorf("expected 2 watchlists, got %d", len(watchlists))
	}
	if _, ok := watchlists["alice"]; !ok {
		t.Errorf("missing watchlist for alice")
	}
}

func TestCollectWatchlists_Error(t *testing.T) {
	provider := &fakeWatchlistProvider{
		results: map[string]*domain.Watchlist{
			"bob": makeWatchlist("bob", "Matrix"),
		},
		errs: map[string]error{
			"alice": errors.New("fetch error"),
		},
	}
	service := NewPickService(provider)

	params := &domain.ScrapperParams{}
	watchlists, err := service.collectWatchlists([]string{"alice", "bob"}, params)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if watchlists == nil {
		t.Fatalf("expected partial results, got nil")
	}
}

func TestPick_CommonFilms(t *testing.T) {
	provider := &fakeWatchlistProvider{
		results: map[string]*domain.Watchlist{
			"alice": makeWatchlist("alice", "Matrix", "Inception"),
			"bob":   makeWatchlist("bob", "Matrix", "Tenet"),
		},
	}
	service := NewPickService(provider)

	params := &domain.PickParams{
		Usernames: []string{"alice", "bob"},
		ScrapperParams: &domain.ScrapperParams{
			Genres:   []string{"action"},
			Platform: "netflix",
		},
		Limit: 10,
	}

	films, err := service.Pick(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []domain.Film{{Title: "Matrix"}}
	if !reflect.DeepEqual(films, expected) {
		t.Errorf("expected %+v, got %+v", expected, films)
	}
}

func TestPick_WithLimit(t *testing.T) {
	provider := &fakeWatchlistProvider{
		results: map[string]*domain.Watchlist{
			"alice": makeWatchlist("alice", "Matrix", "Inception", "Tenet"),
			"bob":   makeWatchlist("bob", "Matrix", "Inception", "Tenet"),
		},
	}
	service := NewPickService(provider)

	params := &domain.PickParams{
		Usernames: []string{"alice", "bob"},
		ScrapperParams: &domain.ScrapperParams{
			Genres: []string{"sci-fi"},
		},
		Limit: 2,
	}

	films, err := service.Pick(params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(films) != 2 {
		t.Errorf("expected 2 films, got %d", len(films))
	}
}
