package repository

import (
	"gotranslate/models"
	"gotranslate/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testData = []models.Resource{
	{Key: "key1", LanguageCode: "en", Text: "text 1"},
	{Key: "key1", LanguageCode: "es", Text: "text 2"},
	{Key: "key1", LanguageCode: "fi", Text: "text 3"},
	{Key: "key2", LanguageCode: "en", Text: "text 4"},
	{Key: "key2", LanguageCode: "es", Text: "text 5"},
	{Key: "key3", LanguageCode: "en", Text: "text 6"},
}

func TestGetResourcesByLanguageCode_ShouldGetCorrectResults(t *testing.T) {
	languageCodeUnderTest, expectedResultsCount := "en", 3
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()
	db.Create(&testData)

	results, err := repo.GetResourcesByLanguageCode(languageCodeUnderTest)

	assert.NoError(t, err)
	assert.Len(t, results, expectedResultsCount)
	for _, r := range results {
		assert.Equal(t, languageCodeUnderTest, r.LanguageCode)
	}
}

func TestGetResourcesByKey_ShouldGetCorrectResults(t *testing.T) {
	keyUnderTest, expectedResultsCount := "key1", 3
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()
	db.Create(&testData)

	results, err := repo.GetResourcesByKey(keyUnderTest)

	assert.NoError(t, err)
	assert.Len(t, results, expectedResultsCount)
	for _, r := range results {
		assert.Equal(t, keyUnderTest, r.Key)
	}
}

func TestExistingLanguageCodes_ShouldReturn3(t *testing.T) {
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()
	db.Create(&testData)

	results, err := repo.ExistingLanguageCodes()

	assert.NoError(t, err)
	assert.Len(t, results, 3)
}

func TestAddResources_WhenManyResources_ShouldSucceedAndResourceShouldBeWrittenInTable(t *testing.T) {
	// arrange
	data := []models.Resource{
		{Key: "testKey1", LanguageCode: "en", Text: "test text1"},
		{Key: "testKey2", LanguageCode: "en", Text: "test text2"},
		{Key: "testKey3", LanguageCode: "en", Text: "test text3"},
		{Key: "testKey4", LanguageCode: "en", Text: "test text4"},
	}
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()

	// act
	err := repo.AddResources(data...)

	// assert
	assert.NoError(t, err)

	var results []models.Resource
	db.Find(&results)
	for _, v := range data {
		assert.Contains(t, results, v)
	}
}

func TestAddResources_When1Resource_ShouldSucceedAndResourceShouldBeWrittenInTable(t *testing.T) {
	// arrange
	expectedKey, expectedLanguageCode, expectedText := "testKey", "en", "test text"
	data := []models.Resource{
		{Key: expectedKey, LanguageCode: expectedLanguageCode, Text: expectedText},
	}
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()

	// act
	err := repo.AddResources(data...)

	// assert
	assert.NoError(t, err)

	var result models.Resource
	db.First(&result)
	assert.Equal(t, expectedKey, result.Key)
	assert.Equal(t, expectedLanguageCode, result.LanguageCode)
	assert.Equal(t, expectedText, result.Text)
}

func TestUpdateResourceValues_ShouldChangeText(t *testing.T) {
	// arrange
	resourceUnderTest := models.Resource{Key: testData[0].Key, LanguageCode: testData[0].LanguageCode, Text: testData[0].Text}
	expectedText, expectedAffectedRows := "something else !", int64(1)
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()
	db.Create(&testData)

	// act
	resourceUnderTest.Text = expectedText
	rowsAffected, err := repo.UpdateResourceValues(resourceUnderTest)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAffectedRows, rowsAffected)

	fetchQuery := db.Model(&models.Resource{}).Where("Key = ? AND LanguageCode = ?", resourceUnderTest.Key, resourceUnderTest.LanguageCode)
	var updatedResource models.Resource
	var updatedTotalResults int64
	fetchQuery.First(&updatedResource)
	fetchQuery.Count(&updatedTotalResults)
	assert.Equal(t, expectedText, updatedResource.Text)
	assert.Equal(t, expectedAffectedRows, updatedTotalResults)
}

func TestRemoveResources_WhenExistsInDb_ShouldBeDeleted(t *testing.T) {
	keyUnderTest, languageCodeUnderTest, expectedAffectedRows := testData[0].Key, testData[0].LanguageCode, int64(1)
	db, teardown := testutils.SpinUpContainer(t)
	defer teardown()
	repo := NewResourceGorm(db)
	repo.Init()
	db.Create(&testData)
	fetchQuery := db.Model(&models.Resource{}).Where("Key = ? AND LanguageCode = ?", keyUnderTest, languageCodeUnderTest)
	var totalResultsFoundBefore, totalResultsFoundAfter int64
	fetchQuery.Count(&totalResultsFoundBefore)

	rowsAffected, err := repo.RemoveResources(keyUnderTest, languageCodeUnderTest)

	assert.Equal(t, totalResultsFoundBefore, int64(1))
	assert.NoError(t, err)
	assert.Equal(t, expectedAffectedRows, rowsAffected)
	fetchQuery.Count(&totalResultsFoundAfter)
	assert.NotEqual(t, totalResultsFoundBefore, totalResultsFoundAfter)
}
