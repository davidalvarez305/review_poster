package handlers

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetTemplates(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var templates []models.Template

	err := database.DB.Where("user_id = ?", userId).Find(&templates).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch templates.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": templates,
	})
}

func CreateTemplate(c *fiber.Ctx) error {
	var template models.Template
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to Parse Request Body.",
		})
	}

	err = database.DB.Where("user_id = ?", userId).Save(&template).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to save templates",
		})
	}

	err = database.DB.Where("user_id = ?", userId).Find(&template).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to find template after saving.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": template,
	})
}

func UpdateTemplate(c *fiber.Ctx) error {
	var template models.Template
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	// Set updateable values aside
	templateName := template.Name

	// Query to find record
	err = database.DB.Where("user_id = ? AND id = ?", userId, template.ID).First(&template).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to find template.",
		})
	}

	// If record is found, update. If not, DB will throw error.
	template.Name = templateName

	err = database.DB.Save(&template).First(&template).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to update template.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": template,
	})
}

func DeleteTemplate(c *fiber.Ctx) error {
	templateId := c.Query("template")
	userId := c.Params("userId")
	var template models.Template

	if len(templateId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in params.",
		})
	}

	err := database.DB.Where("user_id = ? AND id = ?", userId, templateId).Delete(&template).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete template.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
