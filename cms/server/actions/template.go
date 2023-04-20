package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

func GetTemplates(userId string) ([]models.Template, error) {
	var templates []models.Template

	err := database.DB.Where("user_id = ?", userId).Find(&templates).Error

	if err != nil {
		return templates, err
	}

	return templates, nil
}

func GetTemplateByID(id string) (models.Template, error) {
	var template models.Template

	err := database.DB.Where("id = ?", id).Find(&template).Error

	if err != nil {
		return template, err
	}

	return template, nil
}

func CreateTemplate(template models.Template, userId string) error {
	return database.DB.Where("user_id = ?").Save(&template).Error
}

func UpdateTemplate(template models.Template, userId string) error {

	// Set updateable values aside
	templateName := template.Name

	// Query to find record
	err := database.DB.Where("user_id = ? AND id = ?", userId, template.ID).First(&template).Error

	if err != nil {
		return err
	}

	// If record is found, update. If not, DB will throw error.
	template.Name = templateName

	return database.DB.Save(&template).First(&template).Error
}

func DeleteTemplate(template models.Template, userId, template_id string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, template_id).Delete(&template).Error
}
