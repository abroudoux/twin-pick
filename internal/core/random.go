package core

import (
	"fmt"
	"math/rand"
)

func SelectRandomFilm(films []string) (string, error) {
	if len(films) == 0 {
		return "", fmt.Errorf("no films available to select from")
	}

	randNum := rand.Intn(len(films))
	return films[randNum], nil
}
