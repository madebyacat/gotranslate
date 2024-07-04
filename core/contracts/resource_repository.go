package contracts

import "gotranslate/models"

type ResoureRepository interface {
	Init() error
	GetResourcesByLanguageCode(languageCode string) ([]models.Resource, error)
	GetResourcesByKey(key string) ([]models.Resource, error)
	AddResources(resources ...models.Resource) error
	UpdateResourceValues(resources ...models.Resource) (rowsAffected int64, err error)
	RemoveResources(key, languageCode string) (rowsAffected int64, err error)
	ExistingLanguageCodes() ([]models.LanguageResult, error)
}
