package http

import (
	"net/http"
	"strings"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	matchService  *application.MatchService
	commonService *application.CommonService
}

func NewServer(matchService *application.MatchService, commonService *application.CommonService) *Server {
	return &Server{matchService: matchService, commonService: commonService}
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
	v1.GET("/common", s.handleCommon)
}

func (s *Server) handleMatch(c *gin.Context) {
	usernames := strings.Split(c.Query("usernames"), ",")

	rawGenres := strings.Split(c.Query("genres"), ",")
	var genres []string
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform := c.Query("platform")

	params := domain.NewScrapperParams(genres, platform)

	film, err := s.matchService.MatchFilm(usernames, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"selected_film": film.Name})
}

func (s *Server) handleCommon(c *gin.Context) {
	usernames := strings.Split(c.Query("usernames"), ",")

	rawGenres := strings.Split(c.Query("genres"), ",")
	var genres []string
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform := c.Query("platform")

	params := domain.NewScrapperParams(genres, platform)

	films, err := s.commonService.GetCommonFilms(usernames, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"common_films": films})
}
