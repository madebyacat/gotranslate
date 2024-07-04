package rest

import (
	"gotranslate/api/middleware"
	"gotranslate/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(auth models.AuthConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var creds models.Credentials
		if err := ctx.BindJSON(&creds); err != nil {
			badRequest(ctx, "invalid request")
			return
		}

		// Replace this with your user authentication logic
		if creds.Username != auth.Username || creds.Password != auth.Password {
			errorResult(ctx, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := middleware.GenerateToken(creds.Username, auth.JwtKey)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "could not generate token")
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}
