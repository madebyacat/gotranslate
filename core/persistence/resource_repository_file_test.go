package persistence

import (
	"encoding/json"
	"gotranslate/models"
	"os"
	"sync"
	"testing"
)

func TestGetResources(t *testing.T) {
	expected, languageCode := 2, "fi"

	testData := []models.Resource{
		{Key: "mykey", LanguageCode: languageCode, Text: "val1"},
		{Key: "myotherkey", LanguageCode: languageCode, Text: "val2"},
	}
	repo, _, err := createTestfileWithData(true, testData)
	if err != nil {
		t.Error("failed to create the file")
	}

	results, err := repo.GetResourcesByLanguageCode(languageCode)

	if err != nil || len(results) != expected {
		t.Error("test failed")
	}
}

func TestGetResource(t *testing.T) {
	expected, key := 2, "abcde"

	testData := []models.Resource{
		{Key: key, LanguageCode: "no", Text: "val1"},
		{Key: key, LanguageCode: "fi", Text: "val3"},
	}
	repo, _, err := createTestfileWithData(true, testData)
	if err != nil {
		t.Error("failed to create the file")
	}

	results, err := repo.GetResourcesByKey(key)

	if err != nil || len(results) != expected {
		t.Error("test failed")
	}
}

func TestAddResourcesWith1Resource(t *testing.T) {
	// arrange
	expectedText, keyUnderTest := "test value", "myNewKey"
	repo, cleanup, err := createTestfileWithData(true, []models.Resource{})
	if err != nil {
		t.Error("failed to create the file")
	}

	// act
	err = repo.AddResources(models.Resource{Key: keyUnderTest, LanguageCode: "en", Text: expectedText})
	if err != nil {
		t.Error(err.Error())
	}

	data, err := repo.GetResourcesByKey(keyUnderTest)
	if err != nil {
		t.Error(err.Error())
	}

	// assert
	dataCount := len(data)
	if dataCount != 1 {
		t.Errorf("expected 1 result but got %v", dataCount)
	}
	if dataCount > 0 && data[0].Text != expectedText {
		t.Errorf("expected Text %v but got %v, for newly added item with key %v", expectedText, data[0].Text, keyUnderTest)
	}

	cleanup()
}

func TestAddResourcesFile(t *testing.T) {
	// arrange
	keyUnderTest := "myNewKey"
	repo, cleanup, err := createTestfileWithData(true, []models.Resource{})
	if err != nil {
		t.Error("failed to create the file")
	}

	// act
	err = repo.AddResources(models.Resource{Key: keyUnderTest, LanguageCode: "en", Text: "newValue1"},
		models.Resource{Key: keyUnderTest, LanguageCode: "sv", Text: "newValue2"})
	if err != nil {
		t.Error(err.Error())
	}

	data, err := repo.GetResourcesByKey(keyUnderTest)
	if err != nil {
		t.Error(err.Error())
	}

	// assert
	if dataCount := len(data); dataCount != 2 {
		t.Errorf("expected %v items for key %v, but got %v", 2, keyUnderTest, dataCount)
	}

	cleanup()
}

func createTestfileWithData(withFixedData bool, data []models.Resource) (repo *ResourceRepositoryFile, cleanup func(), err error) {
	if withFixedData {
		fixedData := []models.Resource{
			{Key: "aa", LanguageCode: "en", Text: "val1"},
			{Key: "bb", LanguageCode: "en", Text: "val2"},
			{Key: "aa", LanguageCode: "es", Text: "val3"},
			{Key: "bb", LanguageCode: "es", Text: "val4"},
			{Key: "aa", LanguageCode: "sv", Text: "val5"},
		}
		data = append(data, fixedData...)
	}
	repo = &ResourceRepositoryFile{File: "test.txt", Mu: sync.Mutex{}}

	file, err := os.Create(repo.File)
	if err != nil {
		return nil, func() {}, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for _, resource := range data {
		encoder.Encode(resource)
	}

	cleanup = func() {
		os.Remove(repo.File)
	}

	return repo, cleanup, nil
}
