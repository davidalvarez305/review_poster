package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type CreateWordInput struct {
	ID       int      `json:"id"`
	Word     string   `json:"word"`
	Tag      string   `json:"tag"`
	UserID   int      `json:"user_id"`
	Synonyms []string `json:"synonyms"`
}

func GetWordByName(name, userId string) (models.Word, error) {
	var word models.Word

	err := database.DB.Where("name = ? AND user_id = ?", name, userId).Find(&word).Error

	if err != nil {
		return word, err
	}

	return word, nil
}

func GetWords(userId string) ([]models.Word, error) {
	var words []models.Word

	err := database.DB.Where("user_id = ?", userId).Preload("User").Find(&words).Error

	if err != nil {
		return words, err
	}

	return words, nil
}

func CreateWord(word models.Word) error {
	return database.DB.Where("user_id = ?", word.UserID).Save(&word).Error
}

func UpdateWord(word models.Word, userId string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, word.ID).Save(&word).First(&word).Error
}

func DeleteWord(word models.Word, userId, word_id string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, word_id).Delete(&word).Error
}
