package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
)

type Server struct {
	PickService *application.PickService
}

func NewServer(ps *application.PickService) *Server {
	return &Server{PickService: ps}
}

func (s *Server) Run() {
	r := gin.Default()
	s.registerRoutes(r)
	r.Run(":8080")
}

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/pick", s.handlePick)
}

func (s *Server) handlePick(c *gin.Context) {
	usernames := strings.Split(c.Query("usernames"), ",")
	if len(usernames) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "--usernames is required"})
		return
	}

	rawGenres := strings.Split(c.Query("genres"), ",")
	var genres []string
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform := c.Query("platform")

	limit := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	params := domain.NewScrapperParams(genres, platform)

	films, err := s.PickService.Pick(usernames, params, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"films": films})
}
