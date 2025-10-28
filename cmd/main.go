package main

import (
	match "github.com/abroudoux/twinpick/internal/match"
	scrapper "github.com/abroudoux/twinpick/internal/scrapper"
	"github.com/charmbracelet/log"
)

func main() {
	usernames := []string{"abroudoux", "66Sceptre", "mascim"}
	watchlists := scrapper.ScrapUsersWachtlists(usernames)

	for username, films := range watchlists {
		log.Infof("=== %s's Watchlist (%d films) ===", username, len(films))

		for i, film := range films {
			log.Infof("%d: %s", i+1, film)
		}
	}

	commonFilms, err := match.MatchWatchlists(watchlists)
	if err != nil {
		log.Fatal("Error while matching watchlists: ", err)
		return
	}

	for _, film := range commonFilms {
		log.Infof("Common film found: %s", film)
	}

	filmSelected, err := match.SelectRandomFilm(commonFilms)
	if err != nil {
		log.Fatal("Error while selecting a random film: ", err)
		return
	}
	if filmSelected == "" {
		log.Info("No common film found among the watchlists.")
		return
	}

	log.Infof("Film found: %s", filmSelected)
}
