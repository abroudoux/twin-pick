package cli

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "twinpick",
	Short: "Twinpick CLI : find the perfect film based on your Letterboxd Watchlists",
	Long:  "Twinpick is a tool to help you find the perfect film to watch based on your Letterboxd Watchlists.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
