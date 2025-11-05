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

func NewServer(ps application.PickServiceInterface, ss application.SpotServiceInterface) *Server {
	s := &Server{PickService: ps, SpotService: ss}
	s.Router = gin.Default()
	s.registerRoutes(s.Router)
	return s
}

func (s *Server) Run(port string) {
	s.Router.Run(fmt.Sprintf(":%s", port))
}

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/pick", s.handlePick)
	v1.GET("/spot", s.handleSpot)
}
