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
	Short: "Twinpick CLI : find the perfect film based on your Letterboxd Watchlists",
}

var (
	usernames     string
	genres        string
	platform      string
	matchService  *application.MatchService
	commonService *application.CommonService
)

func Execute(m *application.MatchService, c *application.CommonService) {
	matchService = m
	commonService = c

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

	commonCmd := &cobra.Command{
		Use:   "common",
		Short: "Find common movies in users' Letterboxd watchlists",
		RunE:  runCommon,
	}

	matchCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	matchCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	matchCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")

	commonCmd.Flags().StringVar(&usernames, "usernames", "", "Comma-separated Letterboxd usernames (required)")
	commonCmd.Flags().StringVar(&genres, "genres", "", "Optional genres, comma-separated")
	commonCmd.Flags().StringVar(&platform, "platform", "", "Optional platform, e.g., netflix-fr")

	rootCmd.AddCommand(matchCmd)
	rootCmd.AddCommand(commonCmd)
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

	film, err := matchService.MatchFilm(userList, params)
	if err != nil {
		return err
	}

	log.Infof("🎬 Selected film: %s", film.Title)
	return nil
}

func runCommon(cmd *cobra.Command, args []string) error {
	if usernames == "" {
		return fmt.Errorf("--usernames is required")
	}

	userList := strings.Split(usernames, ",")
	genreList := []string{}
	if genres != "" {
		genreList = strings.Split(genres, ",")
	}

	params := domain.NewScrapperParams(genreList, platform)

	_, err := commonService.GetCommonFilms(userList, params)
	if err != nil {
		return err
	}

	return nil
}
