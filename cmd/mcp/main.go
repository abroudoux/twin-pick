package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/mcp"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure"
)

func main() {
	provider := infrastructure.NewLetterboxdScrapper()
	pickService := application.NewPickService(provider)

	server := mcp.NewServer(pickService)
	server.Run()
}
