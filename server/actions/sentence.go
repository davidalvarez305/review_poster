package actions

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
)

func GetUserSentencesByParagraphAndTemplate(paragraph, userId, template string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Preload("Paragraph.Template").Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id").Where("paragraph.name = ? AND template.user_id = ? AND template.name = ?", paragraph, userId, template).Find(&sentences).Error

	if err != nil {
		return sentences, err
	}

	return sentences, nil
}

func GetSentencesByTemplate(template, userId string) ([]models.Sentence, error) {
	var sentences []models.Sentence

	err := database.DB.Preload("Paragraph.Template").Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id").Where("template.user_id = ? AND template.name = ?", userId, template).Find(&sentences)

	if err.Error != nil {
		return sentences, err.Error
	}

	return sentences, nil
}
