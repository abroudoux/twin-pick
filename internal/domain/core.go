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
		Name: name,
	}
}

func NewScrapperParams(genres []string, platform string) *ScrapperParams {
	return &ScrapperParams{
		Genres:   genres,
		Platform: platform,
	}
}

func GetCommonFilms(watchlists map[string]*Watchlist) ([]Film, error) {
	if len(watchlists) == 0 {
		return nil, errors.New("no watchlists provided")
	}

	filmCount := make(map[string]int)
	for _, wl := range watchlists {
		seen := make(map[string]bool)
		for _, f := range wl.Films {
			if !seen[f.Name] {
				filmCount[f.Name]++
				seen[f.Name] = true
			}
		}
	}

	var commonFilms []Film
	for name, count := range filmCount {
		if count == len(watchlists) {
			commonFilms = append(commonFilms, NewFilm(name))
		}
	}

	return commonFilms, nil
}

func SelectRandomFilm(films []Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, errors.New("no films to select from")
	}
	return films[rand.Intn(len(films))], nil
}
