package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/cli"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	provider := scrapper.NewLetterboxdScrapper()
	pickService := application.NewPickService(provider)
	spotService := application.NewSpotService(provider)

	cli.Execute(pickService, spotService)
}
