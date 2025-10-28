package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type server struct {
	router *gin.Engine
}

func NewServer() *server {
	router := gin.Default()
	return &server{router: router}
}

func (s *server) Run() {
	s.registerRoutes()
	s.router.Run(fmt.Sprintf(":%d", 8080))
}

func (s *server) registerRoutes() {
	api := s.router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/match/:usernames", s.handleMatch)
}
