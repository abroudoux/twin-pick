package http

import (
	"net/http"
	"strings"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	matchService *application.MatchService
}

func NewServer(matchService *application.MatchService) *Server {
	return &Server{matchService: matchService}
}

func (s *Server) Run() {
	r := gin.Default()
	s.registerRoutes(r)
	r.Run(":8080")
}

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/match", s.handleMatch)
}

func (s *Server) handleMatch(c *gin.Context) {
	usernames := strings.Split(c.Query("usernames"), ",")
	params := domain.ScrapperParams{
		Genres:   strings.Split(c.Query("genres"), ","),
		Platform: c.Query("platform"),
	}

	film, err := s.matchService.FindCommonFilm(usernames, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"selected_film": film.Name})
}
