package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"gotranslate/core/contracts"
	"gotranslate/models"
	"gotranslate/utils"
	"log"
)

type TranslateLanguageMessage struct {
	Type           string
	SourceLanguage string
	TargetLanguage string
}

var _ contracts.BaseMessage = (*TranslateLanguageMessage)(nil)

func (m *TranslateLanguageMessage) GetType() string {
	return "TranslateLanguage"
}

func (m *TranslateLanguageMessage) SetType() {
	m.Type = m.GetType()
}

type TranslateLanguageHandler struct {
	Repo       contracts.ResoureRepository
	Translator contracts.Translator
}

func (h *TranslateLanguageHandler) HandleMessage(messageBody map[string]interface{}) error {
	data, err := json.Marshal(messageBody)
	if err != nil {
		return errors.New("invalid message body")
	}

	var msg TranslateLanguageMessage
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return errors.New("invalid message type")
	}

	return ConsumeTranslation(&msg, h.Repo, h.Translator)
}

func ConsumeTranslation(msg *TranslateLanguageMessage, repo contracts.ResoureRepository, translator contracts.Translator) error {
	if msg.SourceLanguage == "" || msg.TargetLanguage == "" {
		log.Printf("error processing message %v\n", msg)
		return errors.New("languages not set correctly")
	}

	existingResources, err := repo.GetResourcesByLanguageCode(msg.SourceLanguage)
	if err != nil {
		log.Printf("error loading resources %v\n", err.Error())
	}

	if len(existingResources) == 0 {
		return fmt.Errorf("no resources found for language %v", msg.SourceLanguage)
	}

	var newResources []models.Resource
	batches := utils.SplitToBatches(existingResources, translator.GetBatchLimit())
	for _, batch := range batches {
		newResources, err = translator.TranslateResources(msg.TargetLanguage, batch)
		if err != nil {
			return err
		}
	}

	repo.AddResources(newResources...)

	return nil
}
