package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/http"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure"
)

func main() {
	provider := infrastructure.NewLetterboxdScrapper()
	matchService := application.NewMatchService(provider)
	commonService := application.NewCommonService(provider)

	server := http.NewServer(matchService, commonService)
	server.Run()
}
