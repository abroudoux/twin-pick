package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/http"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	letterboxdScrapper := scrapper.NewLetterboxdScrapper()
	browserScrapper := scrapper.NewBrowserScrapper()

	pickService := application.NewPickService(letterboxdScrapper)
	spotService := application.NewSpotService(browserScrapper)

	server := http.NewServer(pickService, spotService)
	server.Run("8080")
}
