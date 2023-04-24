package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/server"
)

func GetParagraphs(userId string) ([]models.Paragraph, error) {
	var paragraphs []models.Paragraph

	err := server.DB.Where("user_id = ?", userId).Preload("Template").Find(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	return paragraphs, nil
}

func CreateParagraph(paragraph models.Paragraph) error {
	return server.DB.Save(&paragraph).First(&paragraph).Error
}

func CreateParagraphs(paragraphs []models.Paragraph, userId string) error {
	return server.DB.Where("user_id = ?", userId).Save(&paragraphs).Error
}

func UpdateParagraph(paragraph models.Paragraph, userId string) error {
	return server.DB.Where("user_id = ?", userId).Save(&paragraph).First(&paragraph).Error
}

func UpdateParagraphs(paragraphs []models.Paragraph, paragraphId, userId, template string) ([]models.Paragraph, error) {
	err := server.DB.Where("id = ? AND user_id = ?", paragraphId, userId).Save(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	p, err := GetParagraphsByTemplate(template, userId)

	if err != nil {
		return paragraphs, err
	}

	return p, nil
}

func DeleteParagraphs(ids []int, templateId string) ([]models.Paragraph, error) {
	var paragraphs []models.Paragraph

	err := server.DB.Delete(&models.Paragraph{}, ids).Error

	if err != nil {
		return paragraphs, err
	}

	err = server.DB.Where("template_id = ?", templateId).Find(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	return paragraphs, nil
}

func GetParagraphsByTemplate(template, userId string) ([]models.Paragraph, error) {
	var paragraphs []models.Paragraph

	err := server.DB.Where("\"Template\".name = ? AND paragraph.user_id = ?", template, userId).Joins("Template").Find(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	return paragraphs, nil
}

func SimpleDelete(paragraphs []models.Paragraph) error {
	return server.DB.Delete(&paragraphs).Error
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
		err := SimpleDelete(paragraphs)

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkParagraphs(clientParagraphs, existingParagraphs []models.Paragraph, userId string) error {
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
		err := CreateParagraphs(paragraphs, userId)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetTemplatesAndParagraphs(userId string) ([]models.Paragraph, error) {
	var paragraphs []models.Paragraph

	err := server.DB.Where("user_id = ?", userId).Preload("Template").Find(&paragraphs).Error

	if err != nil {
		return paragraphs, err
	}

	return paragraphs, nil
}
