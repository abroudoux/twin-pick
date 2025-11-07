package http

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/application"
)

type Server struct {
	Router      *gin.Engine
	PickService application.PickServiceInterface
	SpotService application.SpotServiceInterface
}

func NewServer(pickService application.PickServiceInterface, spotService application.SpotServiceInterface) *Server {
	server := &Server{PickService: pickService, SpotService: spotService}
	server.Router = gin.Default()
	server.registerRoutes(server.Router)
	return server
}

func (s *Server) Run(port string) {
	s.Router.Run(fmt.Sprintf(":%s", port))
}

func (s *Server) registerRoutes(router *gin.Engine) {
	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/pick", s.handlePick)
	v1.GET("/spot", s.handleSpot)
}
