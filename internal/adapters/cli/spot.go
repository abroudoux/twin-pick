package cli

import (
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/abroudoux/twinpick/internal/domain"
)

func runSpot(cmd *cobra.Command, args []string) error {
	genreList := []string{}
	if genres != "" {
		for _, g := range strings.Split(genres, ",") {
			if trimmed := strings.TrimSpace(g); trimmed != "" {
				genreList = append(genreList, trimmed)
			}
		}
	}

	filters := domain.NewFilters(limit, domain.GetDurationFromString(duration))
	scrapperParams := domain.NewScrapperParams(genreList, platform)
	params := domain.NewParams(filters, scrapperParams)
	spotParams := domain.NewSpotParams(params)

	films, err := spotService.Spot(spotParams)
	if err != nil {
		return err
	}

	if len(films) == 0 {
		log.Infof("No common films found for the given users/filters.")
		return nil
	}

	log.Infof("🎬 Spotted films:")
	for i, f := range films {
		log.Infof("%d. %s", i+1, f.Title)
	}

	return nil
}
