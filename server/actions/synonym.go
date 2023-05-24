package actions

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
)

// Select records from DB by a word string.
func GetUserSynonymsByWord(word, userId string) ([]models.Synonym, error) {
	var synonyms []models.Synonym

	err := database.DB.Preload("Word").Joins("INNER JOIN word ON word.id = synonym.word_id").Where("word.name = ? AND word.user_id = ?", word, userId).Find(&synonyms).Error

	if err != nil {
		return synonyms, err
	}

	return synonyms, nil
}
