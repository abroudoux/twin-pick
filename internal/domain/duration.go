package domain

import "sync"

type Duration int

const (
	Short Duration = iota
	Medium
	Long
)

func (d Duration) String() string {
	switch d {
	case Short:
		return "short"
	case Medium:
		return "medium"
	case Long:
		return "long"
	default:
		return "unknown"
	}
}

func GetDurationFromString(str string) Duration {
	switch str {
	case "short":
		return Short
	case "medium":
		return Medium
	case "long":
		return Long
	default:
		return Long
	}
}

func GetDurationFromInt(i int) Duration {
	switch i {
	case 0:
		return Short
	case 1:
		return Medium
	case 2:
		return Long
	default:
		return Long
	}
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
