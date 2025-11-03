package domain

import (
	"errors"
	"math/rand"
)

func NewWatchlist(username string) *Watchlist {
	return &Watchlist{
		Username: username,
		Films:    []Film{},
	}
}

func NewFilm(name string) Film {
	return Film{
		Title: name,
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

func CompareWatchlists(watchlists map[string]*Watchlist) ([]Film, error) {
	if len(watchlists) == 0 {
		return nil, errors.New("no watchlists provided")
	}

	filmCount := countFilmsAcrossWatchlists(watchlists)
	var common []Film
	for title, count := range filmCount {
		if count == len(watchlists) {
			common = append(common, NewFilm(title))
		}
	}

	return common, nil
}

func countFilmsAcrossWatchlists(watchlists map[string]*Watchlist) map[string]int {
	filmCount := make(map[string]int)
	for _, wl := range watchlists {
		seen := make(map[string]bool)
		for _, f := range wl.Films {
			if !seen[f.Title] {
				filmCount[f.Title]++
				seen[f.Title] = true
			}
		}
	}

	return filmCount
}

func SelectRandomFilm(films []Film) (Film, error) {
	return selectRandomFilmWithRand(films, rand.Intn)
}

func selectRandomFilmWithRand(films []Film, randFn func(int) int) (Film, error) {
	if len(films) == 0 {
		return Film{}, errors.New("no films to select from")
	}
	return films[randFn(len(films))], nil
}
