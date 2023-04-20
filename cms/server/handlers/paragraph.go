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
				"data": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": paragraphs,
		})
	}

	paragraphs, err := actions.GetParagraphs(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	err = actions.CreateParagraphs(paragraphs, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	updatedParagraphs, err := actions.UpdateParagraphs(paragraphs, paragraphId, userId, template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}

func DeleteParagraph(c *fiber.Ctx) error {
	p := c.Query("paragraphs")
	template := c.Query("template")
	ids, err := utils.GetIds(p)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	newParagraphs, err := actions.DeleteParagraphs(ids, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": newParagraphs,
	})
}

func BulkParagraphsUpdate(c *fiber.Ctx) error {
	var paragraphsFromClient []models.Paragraph
	template := c.Query("template")
	userId := c.Params("userId")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in query.",
		})
	}

	err := c.BodyParser(&paragraphsFromClient)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	existingParagraphs, err := actions.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = actions.DeleteBulkParagraphs(paragraphsFromClient, existingParagraphs)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = actions.AddBulkParagraphs(paragraphsFromClient, existingParagraphs, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Re-assign paragraphs to what's now on the database.
	updatedParagraphs, err := actions.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedParagraphs,
	})
}
