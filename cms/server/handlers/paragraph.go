package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetParagraphs(c *fiber.Ctx) error {
	paragraphs := &actions.Paragraphs{}
	template := c.Query("template")
	userId := c.Params("userId")

	// Return paragraphs filtered by template if there's a query.
	if len(template) > 0 {
		err := paragraphs.GetParagraphsByTemplate(template, userId)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": paragraphs,
		})
	}

	err := paragraphs.GetParagraphs(userId)

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
	paragraphs := &actions.Paragraphs{}
	userId := c.Params("userId")

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = paragraphs.CreateParagraphs(userId)

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
	paragraphs := &actions.Paragraphs{}
	userId := c.Params("userId")
	paragraphId := c.Params("paragraphId")
	template := c.Query("template")

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = paragraphs.UpdateParagraphs(paragraphId, userId, template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func DeleteParagraph(c *fiber.Ctx) error {
	paragraphs := &actions.Paragraphs{}
	p := c.Query("paragraphs")
	template := c.Query("template")
	ids, err := utils.GetIds(p)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = paragraphs.DeleteParagraphs(ids, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func BulkParagraphsUpdate(c *fiber.Ctx) error {
	paragraphs := &actions.Paragraphs{}
	template := c.Query("template")
	userId := c.Params("userId")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in query.",
		})
	}

	err := c.BodyParser(&paragraphs)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	existingParagraphs := &actions.Paragraphs{}
	err = existingParagraphs.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.

	var paragraphsToDelete actions.Paragraphs
	var paragraphsToAdd actions.Paragraphs

	paragraphsToDelete.DeleteBulkParagraphs(paragraphs, existingParagraphs)
	paragraphsToAdd.AddBulkParagraphs(paragraphs, existingParagraphs, userId)

	// Re-assign paragraphs to what's now on the database.
	err = paragraphs.GetParagraphsByTemplate(template, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": paragraphs,
	})
}
