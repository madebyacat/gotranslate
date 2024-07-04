package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func okData[T any](ctx *gin.Context, data []T) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func badRequest(ctx *gin.Context, errorMessages ...string) {
	errorResult(ctx, http.StatusBadRequest, errorMessages...)
}

func errorResult(ctx *gin.Context, statusCode int, errorMessages ...string) {
	ctx.JSON(statusCode, gin.H{
		"errors": errorMessages,
	})
}
