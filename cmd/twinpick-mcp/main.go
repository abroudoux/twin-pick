package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/mcp"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	letterboxdScrapper := scrapper.NewLetterboxdScrapper()
	browserScrapper := scrapper.NewBrowserScrapper()

	pickService := application.NewPickService(letterboxdScrapper)
	spotService := application.NewSpotService(browserScrapper)

	server := mcp.NewServer(pickService, spotService)
	server.Run()
}
