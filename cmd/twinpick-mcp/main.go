package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/mcp"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	provider := scrapper.NewLetterboxdScrapper()
	pickService := application.NewPickService(provider)
	spotService := application.NewSpotService(provider)

	server := mcp.NewServer(pickService, spotService)
	server.Run()
}
