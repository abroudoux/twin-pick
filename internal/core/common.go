package core

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func GetCommonFilms(watchlists map[string][]string) ([]string, error) {
	var lists [][]string
	for _, wl := range watchlists {
		lists = append(lists, wl)
	}

	if len(lists) == 0 {
		log.Warn("No watchlists provided")
		return []string{}, fmt.Errorf("no watchlists provided")
	}

	filmCount := make(map[string]string)
	occurences := make(map[string]int)

	for _, wl := range watchlists {
		seen := make(map[string]bool)

		for f := range wl {
			if !seen[wl[f]] {
				occurences[wl[f]]++

				if _, exists := filmCount[wl[f]]; !exists {
					filmCount[wl[f]] = wl[f]
				}

				seen[wl[f]] = true
			}
		}
	}

	var commonFilms []string
	for film, count := range occurences {
		if count == len(watchlists) {
			commonFilms = append(commonFilms, filmCount[film])
		}
	}

	return commonFilms, nil
}
