package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

func GetSentencesByParagraph(paragraph, userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Where("\"Paragraph\".name = ? AND \"Paragraph\".user_id = ?", paragraph, userId).Joins("Paragraph").Joins("Template").Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

func GetAllSentences(userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence
	err := database.DB.Where("user_id = ?", userId).Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

// Create a single sentence. This assumes that input will have been validated elsewhere.
func CreateSentence(sentence models.Sentence) error {
	err := database.DB.Save(&sentence).First(&sentence).Error

	if err != nil {
		return err
	}

	return nil
}

// Create many sentences. This assumes that input will have been validated elsewhere.
func CreateSentences(sentences []models.Sentence, userId string) error {
	return database.DB.Where("user_id = ?", userId).Save(&sentences).Where("user_id = ?", userId).Find(&sentences).Error
}

func UpdateSentences(sentences []models.Sentence, paragraph, userId string) ([]models.Sentence, error) {
	err := database.DB.Where("user_id = ?", userId).Save(&sentences).Error

	if err != nil {
		return sentences, err
	}

	updatedSentences, err := GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return sentences, err
	}

	return updatedSentences, nil
}

func DeleteSentences(ids []int, paragraph, userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Delete(&models.Sentence{}, ids).Error

	if err != nil {
		return sentences, err
	}

	newSentences, err := GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return newSentences, err
	}

	return newSentences, nil
}

// Delete records without checking user or paragraph id. Assumes that this is being checked elsewhere.
func SimpleDeleteSentences(sentences []models.Sentence) error {
	return database.DB.Delete(&sentences).Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func DeleteBulkSentences(clientSentences, existingSentences []models.Sentence) ([]models.Sentence, error) {
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
		err := SimpleDeleteSentences(sentences)

		if err != nil {
			return sentences, err
		}
	}

	return sentences, nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkSentences(clientSentences, existingSentences []models.Sentence, userId string) ([]models.Sentence, error) {
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
		err := CreateSentences(sentences, userId)

		if err != nil {
			return sentences, err
		}
	}

	return sentences, nil
}
