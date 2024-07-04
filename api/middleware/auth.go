package middleware

import (
	"gotranslate/models"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(username string, jwtKey []byte) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &models.Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AuthMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		bearerPrefix := "Bearer "
		tokenStr = strings.TrimPrefix(tokenStr, bearerPrefix)

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}

func AllowAnonymous() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
