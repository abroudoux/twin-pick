package cli

import (
	"fmt"
	"strings"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "twinpick",
	Short: "Twinpick CLI : find the perfect film based on your Letterboxd Watchlists",
}

var (
	usernames string
	genres    string
	platform  string
	service   *application.MatchService
)

func Execute(s *application.MatchService) {
	service = s
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Find a movie based on users' Letterboxd watchlists",
		RunE:  runMatch,
	}

	matchCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	matchCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	matchCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")

	rootCmd.AddCommand(matchCmd)
}

func runMatch(cmd *cobra.Command, args []string) error {
	if usernames == "" {
		return fmt.Errorf("--usernames is required")
	}

	userList := strings.Split(usernames, ",")
	genreList := []string{}
	if genres != "" {
		genreList = strings.Split(genres, ",")
	}

	params := domain.NewScrapperParams(genreList, platform)

	film, err := service.MatchFilm(userList, params)
	if err != nil {
		return err
	}

	log.Infof("🎬 Selected film: %s", film.Name)
	return nil
}
