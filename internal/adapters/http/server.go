package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
)

type Server struct {
	Router      *gin.Engine
	PickService application.PickServiceInterface
}

func NewServer(ps application.PickServiceInterface) *Server {
	s := &Server{PickService: ps}
	s.Router = gin.Default()
	s.registerRoutes(s.Router)
	return s
}

func (s *Server) Run() {
	s.Router.Run(":8080")
}

func (s *Server) registerRoutes(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/pick", s.handlePick)
}

func (s *Server) handlePick(c *gin.Context) {
	params, err := returnPickParams(c)
	if err != nil {
		return
	}

	films, err := s.PickService.Pick(params)
	if err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"films": films})
}

func returnPickParams(c *gin.Context) (*domain.PickParams, error) {
	rawUsernames := c.Query("usernames")
	if strings.TrimSpace(rawUsernames) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "param usernames is required"})
		return nil, fmt.Errorf("usernames is required")
	}
	usernames := strings.Split(rawUsernames, ",")

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

	scrapperParams := domain.NewScrapperParams(genres, platform)
	return domain.NewPickParams(usernames, scrapperParams, limit), nil
}

func respondError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
	c.Abort()
}
