package domain

import (
	"errors"
	"math/rand"
)

func SelectRandomFilm(films []*Film) (*Film, error) {
	return selectRandomFilmWithRand(films, rand.Intn)
}

func selectRandomFilmWithRand(films []*Film, randFn func(int) int) (*Film, error) {
	if len(films) == 0 {
		return nil, errors.New("no films to select from")
	}
	return films[randFn(len(films))], nil
}
