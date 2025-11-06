package cli

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"github.com/abroudoux/twinpick/internal/application"
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
	duration    string
	pickService *application.PickService
	spotService *application.SpotService
)

func Execute(ps *application.PickService, ss *application.SpotService) {
	pickService = ps
	spotService = ss

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
	spotComd := &cobra.Command{
		Use:   "spot",
		Short: "Films suggested based on different criteria",
		RunE:  runSpot,
	}

	pickCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	pickCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	pickCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")
	pickCmd.Flags().IntVar(&limit, "limit", 0, "Limit number of films returned (0 = all)")
	pickCmd.Flags().StringVar(&duration, "duration", "long", "Optional duration filter: short, medium, long")

	spotComd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	spotComd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")
	spotComd.Flags().IntVar(&limit, "limit", 0, "Limit number of films returned (0 = all)")

	rootCmd.AddCommand(pickCmd, spotComd)
}
