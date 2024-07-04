package translators

import (
	"gotranslate/core/contracts"
	"gotranslate/models"

	"github.com/brianvoe/gofakeit/v6"
)

type FakeTranslator struct {
}

var _ contracts.Translator = (*FakeTranslator)(nil)

func (f *FakeTranslator) Translate(tq models.TranslationQuery) (results []string, err error) {
	gofakeit.Seed(0)
	for range tq.Q {
		results = append(results, gofakeit.ProductName())
	}

	return results, nil
}

func (f *FakeTranslator) TranslateResources(targetLanguageCode string, resources []models.Resource) ([]models.Resource, error) {
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

func (f *FakeTranslator) GetBatchLimit() int {
	return 5
}
