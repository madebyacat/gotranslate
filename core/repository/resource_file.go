package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/models"
	"os"
	"sync"
)

// This implementation was mostly to experiment
type ResourceFile struct {
	File string
	Mu   sync.Mutex
}

var _ contracts.ResoureRepository = (*ResourceFile)(nil)

func NewResourceFile(file string) *ResourceFile {
	return &ResourceFile{
		File: file,
		Mu:   sync.Mutex{},
	}
}

// creates file if it doesn't exist, panics if it can't create it
func (repo *ResourceFile) Init() error {
	if _, err := os.Stat(repo.File); os.IsNotExist(err) {
		fmt.Printf("File '%v' does not exist, creating", repo.File)
		file, err := os.Create(repo.File)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fmt.Printf("Created file.")
	} else if err != nil {
		return err
	}

	return nil
}

func (repo *ResourceFile) GetResourcesByLanguageCode(languageCode string) ([]models.Resource, error) {
	return repo.getResources([]resourceFilter{{LanguageCode: languageCode}})
}

func (repo *ResourceFile) GetResourcesByKey(key string) ([]models.Resource, error) {
	return repo.getResources([]resourceFilter{{Key: key}})
}

func (repo *ResourceFile) AddResources(resources ...models.Resource) error {
	if len(resources) == 0 {
		return errors.New("no resources to add")
	}

	var filters []resourceFilter
	for _, resource := range resources {
		filters = append(filters, resourceFilter{Key: resource.Key, LanguageCode: resource.LanguageCode})
	}

	existingResources, err := repo.getResources(filters)
	if err != nil {
		return err
	}
	if existingCount := len(existingResources); existingCount > 0 {
		return fmt.Errorf("%v of the resources you are trying to add already exist", existingCount)
	}

	repo.Mu.Lock()
	defer repo.Mu.Unlock()

	file, err := os.OpenFile(repo.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, resource := range resources {
		jsonData, err := json.Marshal(resource)
		if err != nil {
			return err
		}

		_, err = file.Write(append(jsonData, '\n'))
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *ResourceFile) UpdateResourceValues(resources ...models.Resource) (rowsAffected int64, err error) {
	return 0, errors.New("not implemented")
}

func (repo *ResourceFile) RemoveResources(key, languageCode string) (rowsAffected int64, err error) {
	return 0, errors.New("notimplemented")
}

func (repo *ResourceFile) getResources(filters resourceFilters) ([]models.Resource, error) {
	repo.Mu.Lock()
	defer repo.Mu.Unlock()
	var results []models.Resource = []models.Resource{}

	file, err := os.Open(repo.File)
	if err != nil {
		return results, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		var resource models.Resource
		err = json.Unmarshal(line, &resource)
		if err != nil {
			return results, err
		}

		if filters.Contains(resource.Key, resource.LanguageCode) {
			results = append(results, resource)
		}
	}

	if err = scanner.Err(); err != nil {
		return results, err
	}

	return results, nil
}

func (repo *ResourceFile) ExistingLanguageCodes() (results []models.LanguageResult, err error) {
	panic("unimplemented")
}

type resourceFilters []resourceFilter

func (filters resourceFilters) Contains(key, languageCode string) bool {
	for _, filter := range filters {
		if (filter.LanguageCode == "" || languageCode == filter.LanguageCode) && (filter.Key == "" || key == filter.Key) {
			return true
		}
	}

	return false
}
