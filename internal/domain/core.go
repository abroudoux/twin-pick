package domain

import (
	"errors"
	"math/rand"
	"sync"
)

func NewWatchlist(username string) *Watchlist {
	return &Watchlist{
		Username: username,
		Films:    []*Film{},
	}
}

func NewScrapperParams(genres []string, platform string) *ScrapperParams {
	return &ScrapperParams{
		Genres:   genres,
		Platform: platform,
	}
}

func NewPickParams(usernames []string, scrapperParams *ScrapperParams, limit int, duration Duration) *PickParams {
	return &PickParams{
		Usernames:      usernames,
		ScrapperParams: scrapperParams,
		Limit:          limit,
		Duration:       duration,
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

func FilterFilmsByDuration(films []*Film, duration Duration) []*Film {
	var (
		mu            sync.Mutex
		filteredFilms []*Film
		wg            sync.WaitGroup
	)

	for _, film := range films {
		wg.Add(1)
		go func(film *Film) {
			defer wg.Done()

			var shouldInclude bool
			switch {
			case film.Duration == 0:
				shouldInclude = true
			case duration == Short:
				shouldInclude = film.Duration <= 100
			case duration == Medium:
				shouldInclude = film.Duration <= 120
			case duration == Long:
				shouldInclude = true
			}

			if shouldInclude {
				mu.Lock()
				filteredFilms = append(filteredFilms, film)
				mu.Unlock()
			}
		}(film)
	}

	wg.Wait()
	return filteredFilms
}
