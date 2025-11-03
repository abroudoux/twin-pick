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
	SpotService application.SpotServiceInterface
}

func NewServer(ps application.PickServiceInterface, ss application.SpotServiceInterface) *Server {
	s := &Server{PickService: ps, SpotService: ss}
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
	v1.GET("/spot", s.handleSpot)
}

func (s *Server) handlePick(ctx *gin.Context) {
	params, err := returnPickParams(ctx)
	if err != nil {
		return
	}

	films, err := s.PickService.Pick(params)
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"films": films})
}

func (s *Server) handleSpot(ctx *gin.Context) {
	params, err := returnSpotParams(ctx)
	if err != nil {
		return
	}

	films, err := s.SpotService.Spot(params)
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"films": films})
}

func returnPickParams(ctx *gin.Context) (*domain.PickParams, error) {
	rawUsernames := ctx.Query("usernames")
	if strings.TrimSpace(rawUsernames) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "param usernames is required"})
		return nil, fmt.Errorf("usernames is required")
	}
	usernames := strings.Split(rawUsernames, ",")

	rawGenres := strings.Split(ctx.Query("genres"), ",")
	var genres []string
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform := ctx.Query("platform")

	limit := 0
	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	return domain.NewPickParams(usernames, domain.NewScrapperParams(genres, platform), limit), nil
}

func returnSpotParams(ctx *gin.Context) (*domain.SpotParams, error) {
	rawGenres := strings.Split(ctx.Query("genres"), ",")
	var genres []string
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform := ctx.Query("platform")

	limit := 0
	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	return domain.NewSpotParams(domain.NewScrapperParams(genres, platform), limit), nil
}

func respondError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{"error": msg})
	ctx.Abort()
}
