package domain

type Duration int

const (
	Short Duration = iota
	Medium
	Long
)

const DURATION_SHORT = 100
const DURATION_MEDIUM = 120

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
	if duration == Long {
		return films
	}

	filteredFilms := make([]*Film, 0, len(films))

	for _, film := range films {
		if film.Duration == 0 {
			filteredFilms = append(filteredFilms, film)
			continue
		}

		switch duration {
		case Short:
			if film.Duration <= DURATION_SHORT {
				filteredFilms = append(filteredFilms, film)
			}
		case Medium:
			if film.Duration <= DURATION_MEDIUM {
				filteredFilms = append(filteredFilms, film)
			}
		}
	}

	return filteredFilms
}
