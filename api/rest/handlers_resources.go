package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"gotranslate/core/contracts"
	"gotranslate/models"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetResources(repo contracts.ResoureRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		languageCode, key := ctx.Query("languagecode"), ctx.Query("key")
		err := validateSearchFilter(languageCode, key)
		if err != nil {
			badRequest(ctx, "invalid data")
			return
		}

		var resources []models.Resource
		if len(languageCode) > 0 {
			resources, err = repo.GetResourcesByLanguageCode(languageCode)
		} else if len(key) > 0 {
			resources, err = repo.GetResourcesByKey(key)
		}

		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was a problem retrieving the resources from the database")
			return
		}

		okData(ctx, resources)
	}
}

func AddResources(repo contracts.ResoureRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resources, err := parseResourcesFromRequest(ctx)
		if err != nil {
			log.Println(err)
			badRequest(ctx, "invalid request")
			return
		}

		if errors := validateResourceData(resources); errors.HasErrors() {
			badRequest(ctx, errors.AllErrors()...)
			return
		}

		err = repo.AddResources(resources...)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was an error adding resources to the database")
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{})
	}
}

func parseResourcesFromRequest(ctx *gin.Context) ([]models.Resource, error) {
	resources, empty := []models.Resource{}, []models.Resource{}
	var singleResource models.Resource

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return empty, errors.New("unable to read request body")
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &singleResource); err == nil { // try parsing 1 item before multiple
		resources = append(resources, singleResource)
	} else if err := json.Unmarshal(body, &resources); err != nil { // parse multiple
		return empty, errors.New("invalid data")
	}

	if len(resources) == 0 {
		return empty, errors.New("no resources found")
	}
	return resources, nil
}

func DeleteResources(repo contracts.ResoureRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		languageCode, key := ctx.Query("languagecode"), ctx.Query("key")
		if !languageCodeIsValid(languageCode) || len(key) == 0 {
			badRequest(ctx, "To delete a resource you need to provide a valid languageCode and key")
			return
		}

		_, err := repo.RemoveResources(key, languageCode)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was a problem removing the specified resources")
			return
		}

		ctx.Status(http.StatusNoContent)
	}
}

func UpdateResources(repo contracts.ResoureRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resources, err := parseResourcesFromRequest(ctx)
		if err != nil {
			badRequest(ctx, err.Error())
			return
		}

		if errors := validateResourceData(resources); errors.HasErrors() {
			badRequest(ctx, errors.AllErrors()...)
			return
		}

		rowsAffected, err := repo.UpdateResourceValues(resources...)
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was a problem updating the resources")
			return
		}
		if rowsAffected == 0 {
			errorResult(ctx, http.StatusNotFound, "no items to update found")
		}

		ctx.Status(http.StatusNoContent)
	}
}

func GetAvailableLanguages(repo contracts.ResoureRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		results, err := repo.ExistingLanguageCodes()
		if err != nil {
			errorResult(ctx, http.StatusInternalServerError, "there was a problem retrieving the available languages")
			return
		}

		okData(ctx, results)
	}
}
