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

	var dur domain.Duration
	switch strings.ToLower(strings.TrimSpace(duration)) {
	case "short":
		dur = domain.Short
	case "medium":
		dur = domain.Medium
	case "long", "":
		dur = domain.Long
	default:
		dur = domain.Long
	}

	platform = strings.TrimSpace(platform)

	scrapperParams := domain.NewScrapperParams(genreList, platform)
	pickParams := domain.NewPickParams(userList, scrapperParams, limit, dur)

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
