package scrapper

import (
	"fmt"
	"strings"

	"github.com/abroudoux/twinpick/internal/domain"
)

func buildWatchlistURL(username string, params *domain.ScrapperFilters) string {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist", username)

	if len(params.Genres) > 0 {
		url += "/genre/" + strings.Join(params.Genres, "+")
	}
	if params.Platform != "" {
		url += fmt.Sprintf("/on/%s", params.Platform)
	}
	return url
}

func buildPopularFilmsURL(params *domain.ScrapperFilters) string {
	url := "https://letterboxd.com/films/popular"

	if len(params.Genres) > 0 {
		url += "/genre/" + strings.Join(params.Genres, "+")
	}
	return url
}
