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

func CompareWatchlists(watchlists map[string]*Watchlist, strict bool) ([]*Film, error) {
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
	numWatchlists := len(watchlists)
	majorityThreshold := numWatchlists/2 + 1

	for key, count := range filmCount {
		if count == numWatchlists {
			commonFilms = append(commonFilms, filmMap[key])
		} else if !strict && numWatchlists >= 3 && count >= majorityThreshold {
			commonFilms = append(commonFilms, filmMap[key])
		}
	}

	return commonFilms, nil
}
