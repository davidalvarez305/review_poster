package actions

import (
	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
)

type Template struct {
	*models.Template
}

type Templates []*Template

func (templates *Templates) GetTemplates(userId string) error {
	result := database.DB.Where("user_id = ?", userId).Find(&templates)

	return result.Error
}

func (template *Template) GetTemplateByID(id string) error {
	result := database.DB.Where("id = ?", id).First(&template)

	return result.Error
}

func (template *Template) CreateTemplate() error {
	result := database.DB.Save(&template)

	return result.Error
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

	query = database.DB.Save(&template).First(&template)

	return query.Error
}

func (template *Template) DeleteTemplate(userId, template_id string) error {
	result := database.DB.Where("user_id = ? AND id = ?", userId, template_id).Delete(&template)

	return result.Error
}
