package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
)

var rootCmd = &cobra.Command{
	Use:   "twinpick",
	Short: "Twinpick CLI: pick films based on Letterboxd watchlists",
}

var (
	usernames   string
	genres      string
	platform    string
	limit       int
	pickService *application.PickService
)

func Execute(ps *application.PickService) {
	pickService = ps

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	pickCmd := &cobra.Command{
		Use:   "pick",
		Short: "Pick films from users' watchlists",
		RunE:  runPick,
	}

	pickCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	pickCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	pickCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")
	pickCmd.Flags().IntVar(&limit, "limit", 0, "Limit number of films returned (0 = all)")

	rootCmd.AddCommand(pickCmd)
}

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

	pickParams := domain.NewPickParams(userList, domain.NewScrapperParams(genreList, platform), limit)

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
