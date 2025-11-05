package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/http"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	provider := scrapper.NewLetterboxdScrapper()
	pickService := application.NewPickService(provider)
	spotService := application.NewSpotService(provider)

	server := http.NewServer(pickService, spotService)
	server.Run("8080")
}
