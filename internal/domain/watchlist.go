package domain

import "errors"

type Watchlist struct {
	Username string
	Films    []*Film
}

func NewWatchlist(username string) *Watchlist {
	return &Watchlist{
		Username: username,
		Films:    []*Film{},
	}
}

func CompareWatchlists(watchlists map[string]*Watchlist) ([]*Film, error) {
	if len(watchlists) == 0 {
		return nil, errors.New("no watchlists provided")
	}

	type filmKey struct{ Title string }
	filmMap := make(map[filmKey]*Film)
	filmCount := make(map[filmKey]int)

	for _, wl := range watchlists {
		seen := make(map[filmKey]bool)
		for _, f := range wl.Films {
			key := filmKey{Title: f.Title}
			filmMap[key] = f
			if !seen[key] {
				filmCount[key]++
				seen[key] = true
			}
		}
	}

	var commonFilms []*Film
	for key, count := range filmCount {
		if count == len(watchlists) {
			commonFilms = append(commonFilms, filmMap[key])
		} else if count >= len(watchlists)/2 {
			commonFilms = append(commonFilms, filmMap[key])
		}
	}

	return commonFilms, nil
}
