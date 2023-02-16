package actions

import (
	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
)

type Word struct {
	*models.Word
}

type Words []*Word

type CreateWordInput struct {
	ID       int      `json:"id"`
	Word     string   `json:"word"`
	Tag      string   `json:"tag"`
	UserID   uint     `json:"user_id"`
	Synonyms []string `json:"synonyms"`
}

func (word *Word) GetWordByID(id int) error {
	result := database.DB.Find(&word, "id = ?", id)

	return result.Error
}

func (word *Word) GetWordByName(name, userId string) error {
	result := database.DB.Find(&word, "name = ? AND user_id = ?", name, userId)

	return result.Error
}

func (words *Words) GetWords(userId string) error {
	result := database.DB.Find(&words, "user_id = ?", userId)

	return result.Error
}

func (word *Word) CreateWord() error {
	query := database.DB.Where("user_id = ?", word.UserID).Save(&word)

	return query.Error
}

func (word *Word) UpdateWord(userId string) error {

	// Set updateable values aside
	wordName := word.Name
	wordTag := word.Tag

	// Query to find record
	query := database.DB.Where("user_id = ? AND id = ?", userId, word.ID).First(&word)

	if query.Error != nil {
		return query.Error
	}

	// If record is found, update. If not, DB will throw error.
	word.Name = wordName
	word.Tag = wordTag

	query = database.DB.Save(&word).First(&word)

	return query.Error
}

func (word *Word) DeleteWord(userId, word_id string) error {
	result := database.DB.Where("user_id = ? AND id = ?", userId, word_id).Delete(&word)

	return result.Error
}
