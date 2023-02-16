package actions

import (
	"fmt"

	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
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
func (paragraphs *Paragraphs) CreateParagraphs() error {
	result := database.DB.Save(&paragraphs)
	return result.Error
}

func (paragraph *Paragraph) UpdateParagraph() error {
	result := database.DB.Where("user_id = ?", paragraph.UserID).Save(&paragraph).First(&paragraph)
	return result.Error
}

// Updates multiple records and returns DB values.
func (paragraphs *Paragraphs) UpdateParagraphs(template, userId string) error {
	query := database.DB.Save(&paragraphs)

	if query.Error != nil {
		return query.Error
	}

	err := paragraphs.GetParagraphsByTemplate(template, userId)
	return err
}

func (paragraphs *Paragraphs) DeleteParagraphs(ids []int, template string) error {
	res := database.DB.Delete(&models.Paragraph{}, ids)

	if res.Error != nil {
		return res.Error
	}

	sql := fmt.Sprintf(`SELECT * FROM paragraph WHERE template_id = (
		SELECT id FROM template WHERE name = '%s'
	)`, template)
	result := database.DB.Raw(sql).Scan(&paragraphs)

	return result.Error
}

func (paragraphs *Paragraphs) GetParagraphsByTemplate(template, userId string) error {
	sql := fmt.Sprintf(`SELECT * FROM paragraph WHERE template_id = (
		SELECT id FROM template WHERE name = '%s' AND user_id = '%s'
	)`, template, userId)
	result := database.DB.Raw(sql).Scan(&paragraphs)

	return result.Error
}

func (paragraphs *Paragraphs) SimpleDelete() error {
	res := database.DB.Delete(&paragraphs)
	return res.Error
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
func (paragraphs Paragraphs) AddBulkParagraphs(clientParagraphs, existingParagraphs *Paragraphs) error {
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
		err := paragraphs.CreateParagraphs()

		if err != nil {
			return err
		}
	}

	return nil
}
