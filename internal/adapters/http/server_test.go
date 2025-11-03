package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/domain"
)

type fakePickService struct {
	films []domain.Film
	err   error
}

func (f *fakePickService) Pick(params *domain.PickParams) ([]domain.Film, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.films, nil
}

func TestHandlePick_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeService := &fakePickService{
		films: []domain.Film{
			{Title: "Inception"},
			{Title: "Matrix"},
		},
	}

	server := NewServer(fakeService, nil)
	router := server.Router

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pick?usernames=alice,bob&genres=action,drama&platform=netflix&limit=2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	expectedBody := `{"films":[{"Title":"Inception"},{"Title":"Matrix"}]}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestHandlePick_MissingUsernames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	server := NewServer(nil)
	router := server.Router

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pick", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	expectedBody := `{"error":"param usernames is required"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}

func TestHandlePick_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeService := &fakePickService{
		err: errors.New("watchlist fetch failed"),
	}

	server := NewServer(fakeService)
	router := server.Router

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pick?usernames=alice", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	expectedBody := `{"error":"watchlist fetch failed"}`
	if rec.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rec.Body.String())
	}
}
