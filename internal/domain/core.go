package domain

import (
	"errors"
	"math/rand"
)

func NewWatchlist(username string) *Watchlist {
	return &Watchlist{
		Username: username,
		Films:    []*Film{},
	}
}

func NewFilm(name, detailsEndpoint string) *Film {
	return &Film{
		Title:           name,
		DetailsEndpoint: detailsEndpoint,
	}
}

func NewScrapperParams(genres []string, platform string) *ScrapperParams {
	return &ScrapperParams{
		Genres:   genres,
		Platform: platform,
	}
}

func NewPickParams(usernames []string, scrapperParams *ScrapperParams, limit int) *PickParams {
	return &PickParams{
		Usernames:      usernames,
		ScrapperParams: scrapperParams,
		Limit:          limit,
	}
}

func NewSpotParams(scrapperParams *ScrapperParams, limit int) *SpotParams {
	return &SpotParams{
		ScrapperParams: scrapperParams,
		Limit:          limit,
	}
}

func CompareWatchlists(watchlists map[string]*Watchlist) ([]*Film, error) {
	if len(watchlists) == 0 {
		return nil, errors.New("no watchlists provided")
	}

	type filmKey struct{ Endpoint string }
	filmMap := make(map[filmKey]*Film)
	filmCount := make(map[filmKey]int)

	for _, wl := range watchlists {
		seen := make(map[filmKey]bool)
		for _, f := range wl.Films {
			key := filmKey{Endpoint: f.DetailsEndpoint}
			filmMap[key] = f
			if !seen[key] {
				filmCount[key]++
				seen[key] = true
			}
		}
	}

	var common []*Film
	for key, count := range filmCount {
		if count == len(watchlists) {
			common = append(common, filmMap[key])
		}
	}

	return common, nil
}

func SelectRandomFilm(films []*Film) (*Film, error) {
	return selectRandomFilmWithRand(films, rand.Intn)
}

func selectRandomFilmWithRand(films []*Film, randFn func(int) int) (*Film, error) {
	if len(films) == 0 {
		return nil, errors.New("no films to select from")
	}
	return films[randFn(len(films))], nil
}
