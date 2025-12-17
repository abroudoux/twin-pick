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

func buildPopularFilmsAjaxURL(params *domain.ScrapperFilters) string {
	var url string

	switch params.Order {
	case domain.OrderFilterHighest:
		url = "https://letterboxd.com/films/ajax/by/rating/"
	case domain.OrderFilterNewest:
		url = "https://letterboxd.com/films/ajax/by/release/"
	case domain.OrderFilterShortest:
		url = "https://letterboxd.com/films/ajax/by/shortest/"
	default:
		url = "https://letterboxd.com/films/ajax/popular/"
	}

	if len(params.Genres) > 0 {
		url += "genre/" + strings.Join(params.Genres, "+") + "/"
	}

	if params.Platform != "" {
		url += "on/" + params.Platform + "/"
	}

	return url
}
