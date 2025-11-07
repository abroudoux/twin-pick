package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/cli"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	letterboxdScrapper := scrapper.NewLetterboxdScrapper()
	pickService := application.NewPickService(letterboxdScrapper)
	spotService := application.NewSpotService(letterboxdScrapper)

	cli.Execute(pickService, spotService)
}
