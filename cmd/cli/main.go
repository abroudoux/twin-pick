package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/cli"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure"
)

func main() {
	provider := infrastructure.NewLetterboxdScrapper()
	matchService := application.NewMatchService(provider)
	commonService := application.NewCommonService(provider)

	cli.Execute(matchService, commonService)
}
