package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Sentence struct {
	*models.Sentence
}

type Sentences []*models.Sentence

func (sentences *Sentences) GetSentencesByParagraph(paragraph, userId string) error {
	return database.DB.Where("\"Paragraph\".name = ? AND \"Paragraph\".user_id = ?", paragraph, userId).Joins("Paragraph").Joins("Template").Find(&sentences).Error
}

func (sentences *Sentences) GetAllSentences(userId string) error {
	result := database.DB.Where("user_id = ?", userId).Find(&sentences)
	return result.Error
}

// Create a single sentence. This assumes that input will have been validated elsewhere.
func (sentence *Sentence) CreateSentence() error {
	result := database.DB.Save(&sentence).First(&sentence)
	return result.Error
}

// Create many sentences. This assumes that input will have been validated elsewhere.
func (sentences *Sentences) CreateSentences(userId string) error {
	return database.DB.Where("user_id = ?", userId).Save(&sentences).Where("user_id = ?", userId).Find(&sentences).Error
}

func (sentences *Sentences) UpdateSentences(paragraph, userId string) error {
	err := database.DB.Where("user_id = ?", userId).Save(&sentences).Error

	if err != nil {
		return err
	}

	return sentences.GetSentencesByParagraph(paragraph, userId)
}

func (sentences *Sentences) DeleteSentences(ids []int, paragraph, userId string) error {
	res := database.DB.Delete(&models.Sentence{}, ids)

	if res.Error != nil {
		return res.Error
	}

	err := sentences.GetSentencesByParagraph(paragraph, userId)

	return err
}

// Delete records without checking user or paragraph id. Assumes that this is being checked elsewhere.
func (sentences *Sentences) SimpleDelete() error {
	res := database.DB.Delete(&sentences)
	return res.Error
}

func (p *Paragraphs) GetTemplatesAndParagraphs(userId string) error {
	return database.DB.Where("user_id = ?", userId).Preload("Template").Find(&p).Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func (sentences Sentences) DeleteBulkSentences(clientSentences, existingSentences *Sentences) error {
	for _, existingSentence := range *existingSentences {
		found := false
		for _, clientSentence := range *clientSentences {
			if existingSentence.Sentence == clientSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, existingSentence)
		}
	}

	if len(sentences) > 0 {
		err := sentences.SimpleDelete()

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func (sentences Sentences) AddBulkSentences(clientSentences, existingSentences *Sentences, userId string) error {
	for _, clientSentence := range *clientSentences {
		found := false
		for _, existingSentence := range *existingSentences {
			if clientSentence.Sentence == existingSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, clientSentence)
		}
	}

	if len(sentences) > 0 {
		err := sentences.CreateSentences(userId)

		if err != nil {
			return err
		}
	}

	return nil
}
