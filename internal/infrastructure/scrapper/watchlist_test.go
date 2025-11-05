package scrapper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abroudoux/twinpick/internal/domain"
)

func TestGetTotalPages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<div class="paginate-pages">
				<ul>
					<li class="paginate-page"><a>1</a></li>
					<li class="paginate-page"><a>2</a></li>
					<li class="paginate-page"><a>3</a></li>
				</ul>
			</div>
		`))
	}))
	defer ts.Close()

	s := NewLetterboxdScrapper()

	pages, err := s.GetTotalWatchlistPages(ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pages != 3 {
		t.Errorf("expected 3 pages, got %d", pages)
	}
}

func TestGetFilmsOnPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<div class="poster-grid">
				<li><div class="react-component" data-item-full-display-name="Inception"></div></li>
				<li><div class="react-component" data-item-full-display-name="Matrix"></div></li>
			</div>
		`))
	}))
	defer ts.Close()

	s := NewLetterboxdScrapper()

	films, err := s.GetFilmsOnWatchlistPage(ts.URL, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(films) != 2 {
		t.Errorf("expected 2 films, got %d", len(films))
	}
	if films[0].Title != "Inception" || films[1].Title != "Matrix" {
		t.Errorf("unexpected film titles: %+v", films)
	}
}

func TestGetWatchlist_Fake(t *testing.T) {
	s := NewLetterboxdScrapper()

	s.GetTotalWatchlistPages = func(_ string) (int, error) {
		return 2, nil
	}

	s.GetFilmsOnWatchlistPage = func(_ string, page int) ([]*domain.Film, error) {
		switch page {
		case 1:
			return []*domain.Film{
				{Title: "Film1A"},
				{Title: "Film1B"},
			}, nil
		case 2:
			return []*domain.Film{
				{Title: "Film2A"},
				{Title: "Film2B"},
			}, nil
		default:
			return nil, nil
		}
	}

	params := &domain.ScrapperParams{}
	watchlist, err := s.GetWatchlist("fakeuser", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(watchlist.Films) != 4 {
		t.Errorf("expected 4 films, got %d", len(watchlist.Films))
	}

	expectedTitles := map[string]bool{
		"Film1A": true,
		"Film1B": true,
		"Film2A": true,
		"Film2B": true,
	}

	for _, f := range watchlist.Films {
		if _, ok := expectedTitles[f.Title]; !ok {
			t.Errorf("unexpected film %q found", f.Title)
		} else {
			delete(expectedTitles, f.Title)
		}
	}

	if len(expectedTitles) > 0 {
		t.Errorf("missing expected films: %v", expectedTitles)
	}
}
