package messages

import "gotranslate/core/contracts"

// Using Strategy pattern get the appropriate handler per message.
func GetMessageHandlers(repo contracts.ResoureRepository, translator contracts.Translator) map[string]contracts.MessageHandler {
	return map[string]contracts.MessageHandler{
		(&TranslateLanguageMessage{}).GetType(): &TranslateLanguageHandler{Repo: repo, Translator: translator},
	}
}
