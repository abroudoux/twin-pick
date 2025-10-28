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
