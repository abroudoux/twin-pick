package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/cli"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	letterboxdScrapper := scrapper.NewLetterboxdScrapper()
	browserScrapper := scrapper.NewBrowserScrapper()

	pickService := application.NewPickService(letterboxdScrapper)
	spotService := application.NewSpotService(browserScrapper)

	cli.Execute(pickService, spotService)
}
