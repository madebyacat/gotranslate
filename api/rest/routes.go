package rest

import (
	"gotranslate/api/middleware"
	"gotranslate/core/contracts"
	"gotranslate/models"

	"github.com/gin-gonic/gin"
)

func NewRouter(repo contracts.ResoureRepository, translator contracts.Translator, queueClient contracts.QueueService, auth models.AuthConfig) *gin.Engine {
	r := gin.Default()

	r.POST("/login", Authenticate(auth))

	var authMiddleware gin.HandlerFunc
	if auth.SkipAuthentication {
		authMiddleware = middleware.AllowAnonymous()
	} else {
		authMiddleware = middleware.AuthMiddleware(auth.JwtKey)
	}

	protectedRouter := r.Group("/")
	protectedRouter.Use(authMiddleware)
	{
		protectedRouter.GET("/resources", GetResources(repo))
		protectedRouter.POST("/resources", AddResources(repo))
		protectedRouter.PUT("/resources", UpdateResources(repo))
		protectedRouter.DELETE("/resources", DeleteResources(repo))
		protectedRouter.GET("/resources/languages", GetAvailableLanguages(repo))

		protectedRouter.GET("/translations", TranslateResource(translator))
		protectedRouter.POST("/translations/:sourceLanguageCode/to/:targetLanguageCode", TranslateAllToNewLanguage(repo, translator, queueClient))
	}

	return r
}
