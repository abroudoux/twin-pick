package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleSpot(ctx *gin.Context) {
	params := getSpotParams(ctx)

	films, err := s.SpotService.Spot(params)
	if err != nil {
		respondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"films": films})
}
