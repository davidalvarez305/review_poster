package actions

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
)

func GetParagraphsByTemplate(template, userId string) ([]models.Paragraph, error) {
	var paragraphs []models.Paragraph

	err := database.DB.Where("\"Template\".name = ? AND \"Template\".user_id = ?", template, userId).Joins("Template").Find(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	return paragraphs, nil
}

// Takes structs from the client & deletes them. Does not return records from DB.
func DeleteBulkParagraphs(clientParagraphs, existingParagraphs []models.Paragraph) error {
	var paragraphs []models.Paragraph

	for _, existingParagraph := range existingParagraphs {
		found := false
		for _, clientParagraph := range clientParagraphs {
			if existingParagraph.Name == clientParagraph.Name {
				found = true
			}
		}
		if !found {
			paragraphs = append(paragraphs, existingParagraph)
		}
	}

	if len(paragraphs) > 0 {
		err := database.DB.Delete(&paragraphs).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkParagraphs(clientParagraphs, existingParagraphs []models.Paragraph) error {
	var paragraphs []models.Paragraph

	for _, clientParagraph := range clientParagraphs {
		found := false
		for _, existingParagraph := range existingParagraphs {
			if clientParagraph.Name == existingParagraph.Name {
				found = true
			}
		}
		if !found {
			paragraphs = append(paragraphs, clientParagraph)
		}
	}

	if len(paragraphs) > 0 {
		err := database.DB.Save(&paragraphs).Error

		if err != nil {
			return err
		}
	}

	return nil
}
