package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/abroudoux/twinpick/internal/domain"
)

func runPick(cmd *cobra.Command, args []string) error {
	if usernames == "" {
		return fmt.Errorf("--usernames is required")
	}

	userList := strings.Split(usernames, ",")
	for i := range userList {
		userList[i] = strings.TrimSpace(userList[i])
	}

	var genreList []string
	if genres != "" {
		for _, g := range strings.Split(genres, ",") {
			if trimmed := strings.TrimSpace(g); trimmed != "" {
				genreList = append(genreList, trimmed)
			}
		}
	}

	dur := domain.GetDurationFromString(duration)
	platform = strings.TrimSpace(platform)

	filters := domain.NewFilters(limit, dur)
	scrapperParams := domain.NewScrapperParams(genreList, platform)
	params := domain.NewParams(filters, scrapperParams)
	pickParams := domain.NewPickParams(userList, params)

	log.Infof("▶️ Running pick with usernames=%v, genres=%v, platform=%q, limit=%d, duration=%v",
		userList, genreList, platform, limit, dur.String())

	films, err := pickService.Pick(pickParams)
	if err != nil {
		return fmt.Errorf("failed to pick films: %w", err)
	}

	if len(films) == 0 {
		log.Infof("No common films found for the given users/filters.")
		return nil
	}

	log.Infof("🎬 Picked films:")
	for _, f := range films {
		log.Infof("%s", f.Title)
	}

	return nil
}
