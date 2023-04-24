package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/server"
)

func GetSentences(template, userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := server.DB.Where("\"Template\".user_id = ? AND \"Template\".name = ?", userId, template).Joins("Template").Preload("Paragraph").Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

func GetDictionary(userId string) ([]models.Word, error) {
	var words []models.Word

	err := server.DB.Where("user_id = ?", userId).Preload("Synonyms").Find(&words).Error

	if err != nil {
		return words, err
	}

	return words, nil
}
