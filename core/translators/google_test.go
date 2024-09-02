package translators

import (
	"context"
	"gotranslate/models"
	"gotranslate/slices"
	"testing"

	"cloud.google.com/go/translate"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

type translateClientStub struct{}

func (m *translateClientStub) Translate(ctx context.Context, inputs []string, target language.Tag, opts *translate.Options) ([]translate.Translation, error) {
	resultTexts := []translate.Translation{
		{Text: "translated text 1"},
		{Text: "translated text 2"},
		{Text: "translated text 3"},
	}

	return resultTexts, nil
}

var _ TranslateClient = (*translateClientStub)(nil)

func TestTranslate_ShouldReturnTranslatedResults(t *testing.T) {
	mockGoogleClient := translateClientStub{}
	translator := NewGoogle(&mockGoogleClient)
	input := models.TranslationQuery{Target: "fi", Q: []string{"text 1", "text 2", "text 3"}}

	results, err := translator.Translate(input)

	assert.NoError(t, err)
	assert.NotEmpty(t, results)
	for _, translatedText := range results {
		assert.NotEmpty(t, translatedText)
		assert.NotContains(t, input.Q, translatedText)
	}
}

func TestTranslateResources_ShouldReturnResourcesWithTranslatedText(t *testing.T) {
	expectedLanguageCode, substringFromMock := "fi", "translated text"
	mockGoogleClient := translateClientStub{}
	translator := NewGoogle(&mockGoogleClient)
	inputs := []models.Resource{
		{Key: "key1", LanguageCode: "en", Text: "text 1"},
		{Key: "key2", LanguageCode: "en", Text: "text 2"},
		{Key: "key3", LanguageCode: "en", Text: "text 3"},
	}

	results, err := translator.TranslateResources(expectedLanguageCode, inputs)

	assert.NoError(t, err)
	assert.NotEmpty(t, results)

	assert.True(t, slices.All(results, func(r models.Resource) bool { return r.LanguageCode == expectedLanguageCode }), "all results should have expected LanguageCode")
	for _, translatedResource := range results {
		assert.NotEmpty(t, translatedResource.Text)
		assert.Equal(t, expectedLanguageCode, translatedResource.LanguageCode)
		assert.Contains(t, translatedResource.Text, substringFromMock)
		assert.True(t,
			slices.Contains(inputs, func(r models.Resource) bool { return r.Key == translatedResource.Key }),
			"translated resource's Key should be in input data")
		assert.Equal(t, translatedResource.LanguageCode, expectedLanguageCode)
	}
}
