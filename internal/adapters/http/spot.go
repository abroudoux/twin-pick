package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/domain"
)

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
