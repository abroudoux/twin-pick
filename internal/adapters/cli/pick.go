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
	genreList := []string{}
	if genres != "" {
		for _, g := range strings.Split(genres, ",") {
			if trimmed := strings.TrimSpace(g); trimmed != "" {
				genreList = append(genreList, trimmed)
			}
		}
	}

	pickParams := domain.NewPickParams(userList, domain.NewScrapperParams(genreList, platform), limit, domain.Long)

	films, err := pickService.Pick(pickParams)
	if err != nil {
		return err
	}

	if len(films) == 0 {
		log.Infof("No common films found for the given users/filters.")
		return nil
	}

	log.Infof("🎬 Picked films:")
	for i, f := range films {
		log.Infof("%d. %s", i+1, f.Title)
	}

	return nil
}
