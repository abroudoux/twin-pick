package main

import (
	"github.com/abroudoux/twinpick/internal/api"
)

func main() {
	server := api.NewServer()
	server.Run()
}
