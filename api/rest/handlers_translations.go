package rest

import (
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/core/messages"
	"gotranslate/models"
	"gotranslate/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TranslateResource(translator contracts.Translator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, text, targetLanguage := ctx.Query("key"), ctx.Query("text"), ctx.Query("targetLanguage")
		if key == "" || text == "" || !languageCodeIsValid(targetLanguage) {
			badRequest(ctx, "invalid input")
			return
		}

		tq := models.TranslationQuery{Q: []string{text}, Target: targetLanguage}
		translations, err := translator.Translate(tq)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, fmt.Sprintf("there was a problem with the translation service: %v", err.Error()))
			return
		}

		if translationsCount := len(translations); translationsCount != 1 {
			errorResult(ctx, http.StatusTeapot, fmt.Sprintf("expected 1 result from translation service but found %d", translationsCount))
		}

		results := models.Resource{
			Key:          key,
			LanguageCode: targetLanguage,
			Text:         translations[0],
		}

		okData(ctx, []models.Resource{results})
	}
}

func TranslateAllToNewLanguage(repo contracts.ResoureRepository, translator contracts.Translator, queueClient contracts.QueueService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sourceLanguage, targetLanguage := ctx.Param("sourceLanguageCode"), ctx.Param("targetLanguageCode")

		if !languageCodeIsValid(targetLanguage) || !languageCodeIsValid(sourceLanguage) {
			badRequest(ctx, "language code is invalid")
			return
		}

		existingLanguages, err := repo.ExistingLanguageCodes()
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was a problem retrieving existing languages")
			return
		}

		if utils.Contains(existingLanguages, func(item models.LanguageResult) bool { return item.LanguageCode == targetLanguage }) {
			badRequest(ctx, "target language already exists")
			return
		} else if !utils.Contains(existingLanguages, func(item models.LanguageResult) bool { return item.LanguageCode == sourceLanguage }) {
			badRequest(ctx, "source language doesn't exist")
			return
		}

		message := &messages.TranslateLanguageMessage{
			SourceLanguage: sourceLanguage,
			TargetLanguage: targetLanguage,
		}
		err = queueClient.Publish(message)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "something went wrong while starting the translation")
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}
