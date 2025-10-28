package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/mcp"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure"
)

func main() {
	provider := infrastructure.NewLetterboxdScrapper()
	matchService := application.NewMatchService(provider)
	server := mcp.NewServer(matchService)
	server.Run()
}
