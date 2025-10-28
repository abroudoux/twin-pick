package cli

import (
	"fmt"
	"strings"

	"github.com/abroudoux/twinpick/internal/core"
	"github.com/abroudoux/twinpick/internal/scrapper"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	usernames string
	genres    string
	platform  string
)

func init() {
	matchCmd := &cobra.Command{
		Use:   "match",
		Short: "Find a movie based on Users Letterboxd Watchlists",
		RunE: func(cmd *cobra.Command, args []string) error {
			if usernames == "" {
				return fmt.Errorf("the param --usernames is required")
			}

			uList := strings.Split(usernames, ",")
			gList := []string{}
			if genres != "" {
				gList = strings.Split(genres, ",")
			}

			sParams := scrapper.NewScrapperParams(uList, gList, platform)
			watchlists := scrapper.ScrapUsersWachtlists(sParams)

			commonFilms, err := core.GetCommonFilms(watchlists)
			if err != nil {
				return err
			}

			if len(commonFilms) == 0 {
				log.Info("Any common films found 😢")
				return nil
			}

			selectedFilm, _ := core.SelectRandomFilm(commonFilms)
			log.Infof("🎬 Selected film : %s", selectedFilm)
			return nil
		},
	}

	matchCmd.Flags().StringVar(&usernames, "usernames", "", "Usernames Letterboxd split by comma (required)")
	matchCmd.Flags().StringVar(&genres, "genres", "", "Genres (optional, split by comma)")
	matchCmd.Flags().StringVar(&platform, "platform", "", "Platform (optional, ex: netflix-fr)")

	rootCmd.AddCommand(matchCmd)
}
