package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Template struct {
	*models.Template
}

type Templates []*Template

func (templates *Templates) GetTemplates(userId string) error {
	return database.DB.Where("user_id = ?", userId).Find(&templates).Error
}

func (template *Template) GetTemplateByID(id string) error {
	return database.DB.Where("id = ?", id).First(&template).Error
}

func (template *Template) CreateTemplate(userId string) error {
	return database.DB.Where("user_id = ?").Save(&template).Error
}

func (template *Template) UpdateTemplate(userId string) error {

	// Set updateable values aside
	templateName := template.Name

	// Query to find record
	query := database.DB.Where("user_id = ? AND id = ?", userId, template.ID).First(&template)

	if query.Error != nil {
		return query.Error
	}

	// If record is found, update. If not, DB will throw error.
	template.Name = templateName

	return database.DB.Save(&template).First(&template).Error
}

func (template *Template) DeleteTemplate(userId, template_id string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, template_id).Delete(&template).Error
}
