package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Paragraph struct {
	*models.Paragraph
}

type Paragraphs []*models.Paragraph

func (paragraphs *Paragraphs) GetParagraphs(userId string) error {
	result := database.DB.Where("user_id = ?", userId).Find(&paragraphs)
	return result.Error
}

// Create Single Paragraph
func (paragraph *Paragraph) CreateParagraph() error {
	result := database.DB.Save(&paragraph).First(&paragraph)
	return result.Error
}

// Create Multiple Paragraphs
func (paragraphs *Paragraphs) CreateParagraphs(userId string) error {
	return database.DB.Where("user_id = ?", userId).Save(&paragraphs).Error
}

func (paragraph *Paragraph) UpdateParagraph(userId string) error {
	return database.DB.Where("user_id = ?", userId).Save(&paragraph).First(&paragraph).Error
}

// Updates multiple records and returns DB values.
func (paragraphs *Paragraphs) UpdateParagraphs(paragraphId, userId, templateId string) error {
	err := database.DB.Where("id = ? AND user_id = ?", paragraphId, userId).Save(&paragraphs).Error

	if err != nil {
		return err
	}

	return paragraphs.GetParagraphsByTemplate(templateId, userId)
}

func (paragraphs *Paragraphs) DeleteParagraphs(ids []int, templateId string) error {
	err := database.DB.Delete(&models.Paragraph{}, ids).Error

	if err != nil {
		return err
	}

	return database.DB.Where("template_id = ?", templateId).Find(&paragraphs).Error
}

func (paragraphs *Paragraphs) GetParagraphsByTemplate(templateId, userId string) error {
	return database.DB.Where("template_id = ? AND user_id = ?", templateId, userId).Find(&paragraphs).Error
}

func (paragraphs *Paragraphs) SimpleDelete() error {
	return database.DB.Delete(&paragraphs).Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func (paragraphs Paragraphs) DeleteBulkParagraphs(clientParagraphs, existingParagraphs *Paragraphs) error {
	for _, existingParagraph := range *existingParagraphs {
		found := false
		for _, clientParagraph := range *clientParagraphs {
			if existingParagraph.Name == clientParagraph.Name {
				found = true
			}
		}
		if !found {
			paragraphs = append(paragraphs, existingParagraph)
		}
	}

	if len(paragraphs) > 0 {
		err := paragraphs.SimpleDelete()

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func (paragraphs Paragraphs) AddBulkParagraphs(clientParagraphs, existingParagraphs *Paragraphs, userId string) error {
	for _, clientParagraph := range *clientParagraphs {
		found := false
		for _, existingParagraph := range *existingParagraphs {
			if clientParagraph.Name == existingParagraph.Name {
				found = true
			}
		}
		if !found {
			paragraphs = append(paragraphs, clientParagraph)
		}
	}

	if len(paragraphs) > 0 {
		err := paragraphs.CreateParagraphs(userId)

		if err != nil {
			return err
		}
	}

	return nil
}
