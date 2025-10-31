package main

import (
	"github.com/abroudoux/twinpick/internal/adapters/http"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure"
)

func main() {
	provider := infrastructure.NewLetterboxdScrapper()
	pickService := application.NewPickService(provider)

	server := http.NewServer(pickService)
	server.Run()
}
