package persistence

import (
	"strings"
	"testing"
)

// These tests were made mostly to learn, testing the generated query isn't useful in the real world, testing results is

func TestGenerateQuery_WhenOneFilter_ShouldFilterColumns(t *testing.T) {
	// emulating Theory & InlineData
	tests := []struct {
		name                string
		expectedSubstrings  []string
		filter              resourceFilter
		expectedParamsCount int
	}{
		{"testing for Key", []string{" AND Key = $1"}, resourceFilter{Key: "aaa"}, 1},
		{"testing for LanguageCode", []string{" AND LanguageCode = $1"}, resourceFilter{LanguageCode: "en"}, 1},
		{"testing for Key and LanguageCode",
			[]string{" AND Key = $", " AND LanguageCode = $"},
			resourceFilter{Key: "aaa", LanguageCode: "en"},
			2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params, _ := generateQueryAndParameters(tt.filter)

			for _, expected := range tt.expectedSubstrings {
				if !strings.Contains(query, expected) {
					t.Error("expected filters not found")
				}
			}

			if paramsCount := len(params); paramsCount != tt.expectedParamsCount {
				t.Errorf("expected %d parameter but found %d instead", tt.expectedParamsCount, paramsCount)
			}
		})
	}
}

func TestGenerateQuery_WhenNoFilters_ShouldReturnError(t *testing.T) {
	_, _, err := generateQueryAndParameters()

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestGenerateQuery_WhenManyFilters_ShouldUseJoinToFilter(t *testing.T) {
	expectedQuerySubstrings := []string{"WITH ", "$1::jsonb", "INNER JOIN"}
	filter := []resourceFilter{
		{Key: "aaa"},
		{LanguageCode: "en"},
		{Key: "aaa", LanguageCode: "en"},
	}
	query, params, _ := generateQueryAndParameters(filter...)

	for _, expected := range expectedQuerySubstrings {
		if !strings.Contains(query, expected) {
			t.Errorf("expected to contain %v but not found in query", expected)
		}
	}

	if paramsCount := len(params); paramsCount != 1 {
		t.Errorf("expected 1 parameter but found %d instead", paramsCount)
	}
}
