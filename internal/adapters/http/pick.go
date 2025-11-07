package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) handlePick(ctx *gin.Context) {
	params, err := getPickParams(ctx)
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
