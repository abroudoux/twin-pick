package domain

import (
	"errors"
	"testing"
)

func TestNewWatchlist(t *testing.T) {
	wl := NewWatchlist("alice")

	if wl.Username != "alice" {
		t.Errorf("expected username 'alice', got '%s'", wl.Username)
	}
	if wl.Films == nil {
		t.Error("expected initialized slice for Films, got nil")
	}
	if len(wl.Films) != 0 {
		t.Errorf("expected empty Films slice, got %d", len(wl.Films))
	}
}

func TestNewFilm(t *testing.T) {
	f := NewFilm("Inception")

	if f.Title != "Inception" {
		t.Errorf("expected title 'Inception', got '%s'", f.Title)
	}
}

func TestNewScrapperParams(t *testing.T) {
	p := NewScrapperParams([]string{"action", "drama"}, "netflix-fr")

	if len(p.Genres) != 2 || p.Genres[0] != "action" {
		t.Errorf("expected genres ['action','drama'], got %v", p.Genres)
	}
	if p.Platform != "netflix-fr" {
		t.Errorf("expected platform 'netflix-fr', got '%s'", p.Platform)
	}
}

func TestNewPickParams(t *testing.T) {
	s := NewScrapperParams([]string{"comedy"}, "prime-video")
	p := NewPickParams([]string{"bob", "lucy"}, s, 10)

	if len(p.Usernames) != 2 {
		t.Errorf("expected 2 usernames, got %d", len(p.Usernames))
	}
	if p.ScrapperParams.Platform != "prime-video" {
		t.Errorf("expected platform 'prime-video', got '%s'", p.ScrapperParams.Platform)
	}
	if p.Limit != 10 {
		t.Errorf("expected limit 10, got %d", p.Limit)
	}
}

func TestNewSpotParams(t *testing.T) {
	s := NewScrapperParams([]string{"horror"}, "hulu")
	p := NewSpotParams(s, 5)

	if p.ScrapperParams.Platform != "netflix-fr" {
		t.Errorf("expected platform 'netflix-fr', got '%s'", p.ScrapperParams.Platform)
	}
	if p.Limit != 5 {
		t.Errorf("expected limit 5, got %d", p.Limit)
	}
}

func TestCompareWatchlists_CommonFilms(t *testing.T) {
	w1 := NewWatchlist("alice")
	w1.Films = []Film{{Title: "Inception"}, {Title: "Tenet"}}

	w2 := NewWatchlist("bob")
	w2.Films = []Film{{Title: "Tenet"}, {Title: "Interstellar"}}

	watchlists := map[string]*Watchlist{
		"alice": w1,
		"bob":   w2,
	}

	common, err := CompareWatchlists(watchlists)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(common) != 1 {
		t.Fatalf("expected 1 common film, got %d", len(common))
	}
	if common[0].Title != "Tenet" {
		t.Errorf("expected 'Tenet', got '%s'", common[0].Title)
	}
}

func TestCompareWatchlists_NoCommonFilms(t *testing.T) {
	w1 := NewWatchlist("alice")
	w1.Films = []Film{{Title: "Inception"}}

	w2 := NewWatchlist("bob")
	w2.Films = []Film{{Title: "Matrix"}}

	common, err := CompareWatchlists(map[string]*Watchlist{"a": w1, "b": w2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(common) != 0 {
		t.Errorf("expected no common films, got %d", len(common))
	}
}

func TestCompareWatchlists_NoWatchlists(t *testing.T) {
	_, err := CompareWatchlists(map[string]*Watchlist{})
	if err == nil {
		t.Error("expected error when no watchlists provided, got nil")
	}
}

func TestCountFilmsAcrossWatchlists(t *testing.T) {
	w1 := &Watchlist{Films: []Film{{Title: "A"}, {Title: "B"}, {Title: "A"}}}
	w2 := &Watchlist{Films: []Film{{Title: "A"}, {Title: "C"}}}

	res := countFilmsAcrossWatchlists(map[string]*Watchlist{"w1": w1, "w2": w2})

	if res["A"] != 2 {
		t.Errorf("expected A=2, got %d", res["A"])
	}
	if res["B"] != 1 {
		t.Errorf("expected B=1, got %d", res["B"])
	}
	if res["C"] != 1 {
		t.Errorf("expected C=1, got %d", res["C"])
	}
}

func TestSelectRandomFilm_SingleFilm(t *testing.T) {
	films := []Film{{Title: "Inception"}}
	f, err := SelectRandomFilm(films)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Title != "Inception" {
		t.Errorf("expected 'Inception', got '%s'", f.Title)
	}
}

func TestSelectRandomFilmWithRand(t *testing.T) {
	films := []Film{{Title: "A"}, {Title: "B"}, {Title: "C"}}

	f, err := selectRandomFilmWithRand(films, func(n int) int { return 1 })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if f.Title != "B" {
		t.Errorf("expected film 'B', got '%s'", f.Title)
	}
}

func TestSelectRandomFilmWithRand_NoFilms(t *testing.T) {
	_, err := selectRandomFilmWithRand([]Film{}, func(n int) int { return 0 })
	if !errors.Is(err, errors.New("no films to select from")) && err == nil {
		t.Errorf("expected 'no films to select from' error, got %v", err)
	}
}
