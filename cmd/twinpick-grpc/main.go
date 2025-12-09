package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"

	grpcadapter "github.com/abroudoux/twinpick/internal/adapters/grpc"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/infrastructure/scrapper"
)

func main() {
	letterboxdScrapper := scrapper.NewLetterboxdScrapper()
	pickService := application.NewPickService(letterboxdScrapper)
	spotService := application.NewSpotService(letterboxdScrapper)

	server := grpcadapter.NewServer(pickService, spotService, 50051)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Info("Shutting down gRPC server...")
		server.Stop()
	}()

	if err := server.Start(); err != nil {
		log.Fatal("Failed to start gRPC server", "error", err)
	}
}
