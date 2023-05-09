package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetSentences(c *fiber.Ctx) error {
	userId := c.Params("userId")
	paragraph := c.Query("paragraph")

	if len(paragraph) > 0 {
		sentences, err := actions.GetSentencesByParagraph(paragraph, userId)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": "Failed to fetch sentences by paragraph.",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": sentences,
		})
	}

	var sentences []models.Sentence

	err := database.DB.Where("user_id = ?", userId).Find(&sentences).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch all sentences.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func CreateSentences(c *fiber.Ctx) error {
	var sentences []models.Sentence
	userId := c.Params("userId")

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Where("user_id = ?", userId).Save(&sentences).Where("user_id = ?", userId).Find(&sentences).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create sentences.",
		})
	}

	var paragraphs []models.Paragraph

	err = database.DB.Where("user_id = ?", userId).Preload("Template").Find(&paragraphs).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch templates and paragraphs.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": paragraphs,
	})
}

func UpdateSentences(c *fiber.Ctx) error {
	var sentencesFromClient []models.Sentence
	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	err := c.BodyParser(&sentencesFromClient)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err := database.DB.Where("user_id = ?", userId).Save(&sentencesFromClient).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save sentences.",
		})
	}

	updatedSentences, err := GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph after updating.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": updatedSentences,
	})
}

func DeleteSentence(c *fiber.Ctx) error {
	s := c.Query("sentences")
	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	ids, err := utils.GetIds(s)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to parse ids.",
		})
	}

	err := database.DB.Delete(&models.Sentence{}, ids).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete sentences.",
		})
	}

	newSentences, err := GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": newSentences,
	})
}

func BulkSentencesUpdate(c *fiber.Ctx) error {
	var sentencesFromClient []models.Sentence

	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	if paragraph == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No paragraph in query.",
		})
	}

	err := c.BodyParser(&sentencesFromClient)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	existingSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.

	err = actions.DeleteBulkSentences(sentencesFromClient, existingSentences)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete sentences.",
		})
	}

	err = actions.AddBulkSentences(sentencesFromClient, existingSentences, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to save bulk sentences.",
		})
	}

	// Re-assign sentences to what's now on the database.
	updatedSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}
