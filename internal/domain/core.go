package domain

import (
	"errors"
	"math/rand"
)

func GetCommonFilms(watchlists map[string]Watchlist) ([]Film, error) {
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

	var common []Film
	for name, count := range filmCount {
		if count == len(watchlists) {
			common = append(common, Film{Name: name})
		}
	}

	return common, nil
}

func SelectRandomFilm(films []Film) (Film, error) {
	if len(films) == 0 {
		return Film{}, errors.New("no films to select from")
	}
	return films[rand.Intn(len(films))], nil
}
