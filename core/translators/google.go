package translators

import (
	"context"
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/models"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

type Google struct {
	Client TranslateClient
}

var _ contracts.Translator = (*Google)(nil)

type TranslateClient interface {
	Translate(ctx context.Context, inputs []string, target language.Tag, opts *translate.Options) ([]translate.Translation, error)
}

func NewGoogle(client TranslateClient) contracts.Translator {
	return &Google{Client: client}
}

func (t *Google) Translate(tq models.TranslationQuery) ([]string, error) {
	results, emptyResult := []string{}, []string{}

	lang, err := language.Parse(tq.Target)
	if err != nil {
		return emptyResult, err
	}

	translations, err := t.Client.Translate(context.Background(), tq.Q, lang, nil)
	if err != nil {
		return emptyResult, err
	}

	for _, t := range translations {
		results = append(results, t.Text)
	}

	return results, nil
}

func (t *Google) TranslateResources(targetLanguageCode string, resources []models.Resource) ([]models.Resource, error) {
	results, emptyResult := []models.Resource{}, []models.Resource{}
	if len(resources) == 0 {
		return emptyResult, nil
	}

	if len(resources) > t.GetBatchLimit() {
		return emptyResult, fmt.Errorf("there's a limit of processing %d resources per request but found %d", t.GetBatchLimit(), len(resources))
	}

	var textsToTranslate []string
	for _, resource := range resources {
		textsToTranslate = append(textsToTranslate, resource.Text)
	}

	query := models.TranslationQuery{
		Q:      textsToTranslate,
		Target: targetLanguageCode,
	}

	translations, err := t.Translate(query)
	if err != nil {
		return emptyResult, err
	}

	if translationsCount, resourcesCount := len(translations), len(resources); translationsCount != resourcesCount {
		return emptyResult, fmt.Errorf("translation was done but expected %d results and got %d", resourcesCount, translationsCount)
	}

	for i, text := range translations {
		results = append(results, models.Resource{Key: resources[i].Key, LanguageCode: targetLanguageCode, Text: text})
	}

	return results, nil
}

func (f *Google) GetBatchLimit() int {
	return 100 // the actual limit is 128
}
