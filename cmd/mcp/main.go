package main

import "github.com/abroudoux/twinpick/internal/mcp"

func main() {
	server := mcp.NewServer()
	server.Run()
}
