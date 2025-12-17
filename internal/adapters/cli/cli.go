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
	strict      bool
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

	initFlags(pickCmd, spotComd)

	rootCmd.AddCommand(pickCmd, spotComd)
}
