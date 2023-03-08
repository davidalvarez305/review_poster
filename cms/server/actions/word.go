package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Word struct {
	*models.Word
}

type Words []*Word

type CreateWordInput struct {
	ID       int      `json:"id"`
	Word     string   `json:"word"`
	Tag      string   `json:"tag"`
	UserID   int      `json:"user_id"`
	Synonyms []string `json:"synonyms"`
}

func (word *Word) GetWordByName(name, userId string) error {
	return database.DB.Where("name = ? AND user_id = ?", name, userId).Find(&word).Error
}

func (words *Words) GetWords(userId string) error {
	return database.DB.Where("user_id = ?", userId).Preload("User").Find(&words).Error
}

func (word *Word) CreateWord() error {
	return database.DB.Where("user_id = ?", word.UserID).Save(&word).Error
}

func (word *Word) UpdateWord(userId string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, word.ID).Save(&word).First(&word).Error
}

func (word *Word) DeleteWord(userId, word_id string) error {
	result := database.DB.Where("user_id = ? AND id = ?", userId, word_id).Delete(&word)

	return result.Error
}
