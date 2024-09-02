package translators

import (
	"gotranslate/core/contracts"
	"gotranslate/models"

	"github.com/brianvoe/gofakeit/v6"
)

type Fake struct {
}

var _ contracts.Translator = (*Fake)(nil)

func (f *Fake) Translate(tq models.TranslationQuery) (results []string, err error) {
	gofakeit.Seed(0)
	for range tq.Q {
		results = append(results, gofakeit.ProductName())
	}

	return results, nil
}

func (f *Fake) TranslateResources(targetLanguageCode string, resources []models.Resource) ([]models.Resource, error) {
	results := []models.Resource{}
	gofakeit.Seed(0)
	for _, r := range resources {
		newResource := models.Resource{
			Key:          r.Key,
			Text:         gofakeit.ProductName(),
			LanguageCode: targetLanguageCode,
		}
		results = append(results, newResource)
	}

	return results, nil
}

func (f *Fake) GetBatchLimit() int {
	return 5
}
