package rest

import (
	"errors"
	"fmt"
	"gotranslate/models"
)

type ValidationErrors struct {
	Errors []error
}

func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation errors occurred: %d errors", len(ve.Errors))
}

func (ve *ValidationErrors) Add(err error) {
	ve.Errors = append(ve.Errors, err)
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationErrors) AllErrors() []string {
	var results []string
	for _, err := range ve.Errors {
		results = append(results, err.Error())
	}
	return results
}

func validateSearchFilter(languageCode string, key string) error {
	languageCodeSelected, keySelected := len(languageCode) > 0, len(key) > 0

	if !languageCodeSelected && !keySelected {
		return errors.New("filter by 'languagecode' or 'key'")
	}

	if languageCodeSelected && keySelected {
		return errors.New("you can't use both filters")
	}

	if languageCodeSelected && !languageCodeIsValid(languageCode) {
		return errors.New("language code must be 2 letters")
	}

	return nil
}

func validateResourceData(data []models.Resource) *ValidationErrors {
	var validationErrors ValidationErrors

	for _, item := range data {
		if item.LanguageCode != "" && !languageCodeIsValid(item.LanguageCode) {
			validationErrors.Add(errors.New("language code must be 2 letters"))
		}

		if item.Key == "" {
			validationErrors.Add(errors.New("invalid key"))
		}

		if item.Text == "" {
			validationErrors.Add(errors.New("invalid text"))
		}
	}

	return &validationErrors
}

func languageCodeIsValid(languageCode string) bool {
	return len(languageCode) == 2
}
