package messages

import (
	"gotranslate/core/repository"
	"gotranslate/core/translators"
	"gotranslate/models"
	"gotranslate/slices"
	"gotranslate/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsumeTranslation_ShouldSaveTranslatedResourcesToDb(t *testing.T) {
	expectedLanguage := "es"
	data := []models.Resource{
		{Key: "key1", LanguageCode: "en", Text: "text 1"},
		{Key: "key1", LanguageCode: "de", Text: "text 2"},
		{Key: "key1", LanguageCode: "fi", Text: "text 3"},
		{Key: "key2", LanguageCode: "en", Text: "text 4"},
		{Key: "key2", LanguageCode: "de", Text: "text 5"},
		{Key: "key3", LanguageCode: "en", Text: "text 6"},
	}
	msg := TranslateLanguageMessage{Type: "testType", SourceLanguage: "en", TargetLanguage: expectedLanguage}
	translator := translators.FakeTranslator{}
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := repository.NewResourceGorm(db)
	repo.Init()
	db.Create(data)

	ConsumeTranslation(&msg, repo, &translator)

	results, err := repo.GetResourcesByLanguageCode("es")
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	allResultsAreEs := slices.All(results, func(x models.Resource) bool { return x.LanguageCode == expectedLanguage })
	assert.Truef(t, allResultsAreEs, "expected all results to have target language %v", expectedLanguage)
}
