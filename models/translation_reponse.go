package models

type TranslationResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}
