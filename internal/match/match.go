package match

import (
	"github.com/charmbracelet/log"
)

func MatchWatchlists(watchlists map[string][]string) ([]string, error) {
	var lists [][]string
	for _, wl := range watchlists {
		lists = append(lists, wl)
	}

	commonFilms, err := getCommonFilms(lists)
	if err != nil {
		log.Error("Error while getting common films: ", err)
		return []string{}, err
	}

	return commonFilms, nil
}
