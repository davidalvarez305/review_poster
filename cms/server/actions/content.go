package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Content []*models.Sentence

type Dictionary []*models.Word

func (c *Content) GetSentences(template, userId string) error {
	return database.DB.Where("user_id = ? AND name = ?", userId, template).Joins("Template").Preload("Paragraph").Find(&c).Error
}

func (d *Dictionary) GetDictionary(userId string) error {
	return database.DB.Where("user_id = ?", userId).Preload("Word.Synonyms").Find(&d).Error
}
