package contracts

import "gotranslate/models"

type Translator interface {
	Translate(tq models.TranslationQuery) ([]string, error)
	TranslateResources(targetLanguageCode string, resources []models.Resource) ([]models.Resource, error)
	GetBatchLimit() int
}
