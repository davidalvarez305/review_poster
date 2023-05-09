package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetParagraphs(c *fiber.Ctx) error {
	template := c.Query("template")
	userId := c.Params("userId")

	// Return paragraphs filtered by template if there's a query.
	if len(template) > 0 {
		paragraphs, err := actions.GetParagraphsByTemplate(template, userId)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": "Failed to fetch paragraphs.",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": paragraphs,
		})
	}

	// Return all paragraphs without filter
	var paragraphs []models.Paragraph

	err := database.DB.Where("user_id = ?", userId).Preload("Template").Find(&paragraphs).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs.",
		})
	}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func CreateParagraphs(c *fiber.Ctx) error {
	var paragraphs []models.Paragraph

	userId := c.Params("userId")

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Where("user_id = ?", userId).Save(&paragraphs).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to create paragraphs.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func UpdateParagraphs(c *fiber.Ctx) error {
	var paragraphs []models.Paragraph

	userId := c.Params("userId")
	paragraphId := c.Params("paragraphId")
	template := c.Query("template")

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Where("id = ? AND user_id = ?", paragraphId, userId).Save(&paragraphs).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to update paragraphs.",
		})
	}

	err := database.DB.Where("\"Template\".name = ? AND paragraph.user_id = ?", template, userId).Joins("Template").Find(&paragraphs).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template after updating.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}

func DeleteParagraph(c *fiber.Ctx) error {
	paragraphsToDelete := c.Query("paragraphs")
	template := c.Query("template")

	if len(template) == 0 || len(paragraphsToDelete) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template or paragraphs found in querystring.",
		})
	}

	ids, err := utils.GetIds(paragraphsToDelete)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to parse paragraph ID's.",
		})
	}

	var paragraphs []models.Paragraph

	err := database.DB.Delete(&models.Paragraph{}, ids).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete paragraphs.",
		})
	}

	err = database.DB.Where("template_id = ?", template).Find(&paragraphs).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs after deletion.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func BulkParagraphsUpdate(c *fiber.Ctx) error {
	var paragraphsFromClient []models.Paragraph
	template := c.Query("template")
	userId := c.Params("userId")

	if len(template) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in querystring.",
		})
	}

	err := c.BodyParser(&paragraphsFromClient)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	existingParagraphs, err := actions.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template.",
		})
	}

	err = actions.DeleteBulkParagraphs(paragraphsFromClient, existingParagraphs)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete paragraphs in bulk.",
		})
	}

	err = actions.AddBulkParagraphs(paragraphsFromClient, existingParagraphs, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to add new bulk paragraphs.",
		})
	}

	// Re-assign paragraphs to what's now on the database.
	updatedParagraphs, err := actions.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch paragraphs by template after saving.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}
