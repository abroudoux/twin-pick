package domain

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Film struct {
	Title           string   `json:"-"`
	Duration        int      `json:"-"`
	Directors       []string `json:"-"`
	DetailsEndpoint string   `json:"-"`
	Year            int      `json:"-"`
}

func NewFilm(name, detailsEndpoint string) *Film {
	return &Film{
		Title:           name,
		DetailsEndpoint: detailsEndpoint,
	}
}

func (f *Film) MarshalJSON() ([]byte, error) {
	type Alias Film

	if f.Year == 0 {
		if yearFromTitle := extractYearFromTitle(f.Title); yearFromTitle != 0 {
			f.Year = yearFromTitle
		}
	}

	return json.Marshal(&struct {
		*Alias
		Title     string `json:"title"`
		Duration  string `json:"duration"`
		Directors string `json:"directors"`
		Year      string `json:"year"`
	}{
		Alias:     (*Alias)(f),
		Title:     cleanTitle(f.Title),
		Duration:  formatDuration(f.Duration),
		Directors: formatDirectors(f.Directors),
		Year:      formatYear(f.Year),
	})
}

func cleanTitle(title string) string {
	re := regexp.MustCompile(`\s*\(\d{4}\)$`)
	return strings.TrimSpace(re.ReplaceAllString(title, ""))
}

func extractYearFromTitle(title string) int {
	re := regexp.MustCompile(`\((\d{4})\)`)
	matches := re.FindStringSubmatch(title)

	if len(matches) > 1 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			return year
		}
	}
	return 0
}

func formatDuration(duration int) string {
	if duration == 0 {
		return "not found"
	}
	return fmt.Sprintf("%d min", duration)
}

func formatDirectors(directors []string) string {
	if len(directors) == 0 {
		return "not found"
	}
	return strings.Join(directors, ", ")
}

func formatYear(year int) string {
	if year == 0 {
		return "not found"
	}
	return fmt.Sprintf("%d", year)
}
