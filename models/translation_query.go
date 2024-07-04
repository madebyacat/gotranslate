package models

type TranslationQuery struct {
	Q      []string `json:"q"`
	Target string   `json:"target"`
}

func (tq *TranslationQuery) AddText(text string) {
	tq.Q = append(tq.Q, text)
}
