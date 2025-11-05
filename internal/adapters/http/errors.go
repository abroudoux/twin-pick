package http

import "github.com/gin-gonic/gin"

func respondError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{"error": msg})
	ctx.Abort()
}
