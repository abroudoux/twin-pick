package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/abroudoux/twinpick/internal/domain"
)

func parseParams(ctx *gin.Context) (params *domain.Params) {
	var genres []string
	var platform string
	var order domain.OrderFilter

	rawGenres := strings.Split(ctx.Query("genres"), ",")
	for _, g := range rawGenres {
		if trimmed := strings.TrimSpace(g); trimmed != "" {
			genres = append(genres, trimmed)
		}
	}

	platform = ctx.Query("platform")

	order = domain.OrderFilterPopular
	if o := ctx.Query("order"); o != "" {
		switch strings.ToLower(o) {
		case "popular":
			order = domain.OrderFilterPopular
		case "rating", "highest-rated":
			order = domain.OrderFilterHighest
		case "newest":
			order = domain.OrderFilterNewest
		case "shortest":
			order = domain.OrderFilterShortest
		}
	}

	filters := parseFiltsers(ctx)
	scrapperFilters := domain.NewScrapperFilters(genres, platform, order)

	return domain.NewParams(filters, scrapperFilters)
}

func parseFiltsers(ctx *gin.Context) (filters *domain.Filters) {
	var limit int
	var duration domain.Duration

	limit = 0
	if l := ctx.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}

	duration = domain.Long
	if d := ctx.Query("duration"); d != "" {
		switch strings.ToLower(d) {
		case "short":
			duration = domain.Short
		case "medium":
			duration = domain.Medium
		case "long":
			duration = domain.Long
		}
	}

	return domain.NewFilters(limit, duration)
}

func getPickParams(ctx *gin.Context) (pickParams *domain.PickParams, err error) {
	rawUsernames := ctx.Query("usernames")
	if strings.TrimSpace(rawUsernames) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "param usernames is required"})
		return nil, fmt.Errorf("usernames is required")
	}
	usernames := strings.Split(rawUsernames, ",")

	params := parseParams(ctx)

	return domain.NewPickParams(usernames, params), nil
}

func getSpotParams(ctx *gin.Context) (spotParams *domain.SpotParams) {
	params := parseParams(ctx)

	return domain.NewSpotParams(params)
}
