package models

type Resource struct {
	Key          string `gorm:"column:key"`
	LanguageCode string `gorm:"column:languagecode"`
	Text         string `gorm:"column:text"`
}

func (Resource) TableName() string {
	return "resources"
}
