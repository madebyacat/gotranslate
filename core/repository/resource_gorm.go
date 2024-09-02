package repository

import (
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/models"

	"gorm.io/gorm"
)

type ResourceGorm struct {
	DB *gorm.DB
}

func NewResourceGorm(db *gorm.DB) *ResourceGorm {
	return &ResourceGorm{DB: db}
}

var _ contracts.ResoureRepository = (*ResourceGorm)(nil)

func (repo *ResourceGorm) Init() error {
	repo.DB.AutoMigrate(&models.Resource{})
	return nil
}

func (repo *ResourceGorm) AddResources(resources ...models.Resource) error {
	result := repo.DB.CreateInBatches(resources, 10)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ResourceGorm) GetResourcesByKey(key string) ([]models.Resource, error) {
	var resources []models.Resource
	result := repo.DB.Where("Key = ?", key).Find(&resources)
	if result.Error != nil {
		return []models.Resource{}, result.Error
	}

	return resources, nil
}

func (repo *ResourceGorm) GetResourcesByLanguageCode(languageCode string) ([]models.Resource, error) {
	var resources []models.Resource
	result := repo.DB.Where("LanguageCode = ?", languageCode).Find(&resources)
	if result.Error != nil {
		return []models.Resource{}, result.Error
	}

	return resources, nil
}

func (repo *ResourceGorm) RemoveResources(key string, languageCode string) (rowsAffected int64, err error) {
	result := repo.DB.Where("Key = ? AND LanguageCode = ?", key, languageCode).Delete(&models.Resource{})
	if result.Error != nil {
		return 0, result.Error
	}
	fmt.Println(result.RowsAffected)
	return result.RowsAffected, nil
}

func (repo *ResourceGorm) UpdateResourceValues(resources ...models.Resource) (rowsAffected int64, err error) {
	rowsAffected = 0
	for _, resource := range resources {
		result := repo.DB.Model(&models.Resource{}).
			Where("Key = ? AND LanguageCode = ?", resource.Key, resource.LanguageCode).
			Update("Text", resource.Text)
		if result.Error != nil {
			return rowsAffected, result.Error
		}
		rowsAffected += result.RowsAffected
	}

	return rowsAffected, nil
}

func (repo *ResourceGorm) ExistingLanguageCodes() (results []models.LanguageResult, err error) {
	queryResult := repo.DB.
		Model(&models.Resource{}).
		Select(`languagecode as "LanguageCode", COUNT(*) as "Count"`).
		Group("languagecode").
		Scan(&results)
	if queryResult.Error != nil {
		return nil, queryResult.Error
	}

	return results, nil
}
