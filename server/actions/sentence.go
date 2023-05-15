package actions

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
)

func GetSentencesByParagraph(paragraph, userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Where("\"Paragraph\".name = ? AND \"Template\".user_id = ?", paragraph, userId).Joins("Paragraph").Joins("Template").Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

func GetSentencesByTemplate(template string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Where("\"Template\".name = ?", template).Joins("Template").Preload("Paragraph").Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

// Takes structs from the client & deletes them. Does not return records from DB.
func DeleteBulkSentences(clientSentences, existingSentences []models.Sentence) error {
	var sentences []models.Sentence

	for _, existingSentence := range existingSentences {
		found := false
		for _, clientSentence := range clientSentences {
			if existingSentence.Sentence == clientSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, existingSentence)
		}
	}

	if len(sentences) > 0 {
		err := database.DB.Delete(&sentences).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkSentences(clientSentences, existingSentences []models.Sentence) error {
	var sentences []models.Sentence

	for _, clientSentence := range clientSentences {
		found := false
		for _, existingSentence := range existingSentences {
			if clientSentence.Sentence == existingSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, clientSentence)
		}
	}

	if len(sentences) > 0 {
		err := database.DB.Save(&sentences).Error

		if err != nil {
			return err
		}
	}

	return nil
}
