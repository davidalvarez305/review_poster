package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
	"github.com/gofiber/fiber/v2"
)

func GetUserParagraphsByTemplate(c *fiber.Ctx) error {
	template := c.Params("templateName")
	userId := c.Params("userId")

	// Return paragraphs filtered by template if there's a query.
	if len(template) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in URL params.",
		})
	}

	paragraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func CreateUserParagraphsByTemplate(c *fiber.Ctx) error {
	template := c.Params("templateName")
	userId := c.Params("userId")

	if len(template) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in URL params.",
		})
	}

	var paragraphs []models.Paragraph

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&paragraphs).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to save paragraphs.",
		})
	}

	createdParagraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch user's paragraphs by template.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": createdParagraphs,
	})
}

func UpdateUserParagraphsByTemplate(c *fiber.Ctx) error {
	template := c.Params("templateName")
	userId := c.Params("userId")

	if len(template) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in URL params.",
		})
	}

	var input types.UpdateUserParagraphsByTemplateInput

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	if len(input.Paragraphs) > 0 {
		err = database.DB.Save(&input.Paragraphs).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to save paragraphs.",
			})
		}
	}

	// Re-assign paragraphs to what's now on the database.
	updatedParagraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}

func DeleteUserParagraphsByTemplate(c *fiber.Ctx) error {
	template := c.Params("templateName")
	userId := c.Params("userId")

	if len(template) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in URL params.",
		})
	}

	var input types.DeleteUserParagraphsByTemplateInput

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	if len(input.DeleteParagraphs) > 0 {
		err = database.DB.Delete(&models.Paragraph{}, input.DeleteParagraphs).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to delete paragraphs.",
			})
		}
	}

	// Re-assign paragraphs to what's now on the database.
	updatedParagraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}

func UpdateUserParagraphByTemplate(c *fiber.Ctx) error {
	var paragraph models.Paragraph
	template := c.Params("templateName")
	paragraphId := c.Params("paragraphId")

	if len(template) == 0 || len(paragraphId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in URL params.",
		})
	}

	err := c.BodyParser(&paragraph)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user ID from session.",
		})
	}

	var existingParagraph models.Paragraph

	err = database.DB.Joins("JOIN template ON template.id = paragraph.template_id").Where("paragraph.id = ? AND template.user_id = ? AND template.name = ?", paragraph.ID, userId, template).Find(&existingParagraph).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Paragraph does not exist in database.",
		})
	}

	existingParagraph.Name = paragraph.Name
	existingParagraph.Order = paragraph.Order
	existingParagraph.TemplateID = paragraph.TemplateID

	err = database.DB.Save(&existingParagraph).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save paragraphs.",
		})
	}

	updatedParagraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}

func DeleteUserParagraphByTemplate(c *fiber.Ctx) error {
	template := c.Params("templateName")
	userId := c.Params("userId")
	paragraphId := c.Params("paragraphId")

	if len(template) == 0 || len(paragraphId) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"data": "Template or paragraph not in URL params.",
		})
	}

	var paragraph models.Paragraph

	err := database.DB.Joins("INNER JOIN template ON template.id = paragraph.template_id").Where("template.name = ? AND template.user_id = ? AND paragraph.id = ?", template, userId, paragraphId).First(&paragraph).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to find paragraph.",
		})
	}

	err = database.DB.Delete(&paragraph).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete paragraph.",
		})
	}

	paragraphs, err := actions.GetUserParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": paragraphs,
	})
}
