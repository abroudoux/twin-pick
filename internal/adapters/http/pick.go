package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/domain"
)

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

	duration := domain.Long
	if d := ctx.Query("duration"); d != "" {
		switch d {
		case "short":
			duration = domain.Short
		case "medium":
			duration = domain.Medium
		}
	}

	return domain.NewPickParams(usernames, domain.NewScrapperParams(genres, platform), limit, duration), nil
}
