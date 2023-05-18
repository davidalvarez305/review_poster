package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/utils"
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

	err := database.DB.Preload("Paragraph.Template.User").Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id INNER JOIN \"user\" ON \"user\".id = template.user_id").Where("\"user\".id = ?", userId).Find(&sentences).Error

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

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&sentences).Find(&sentences).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create sentences.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": "OK!",
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

	err = database.DB.Save(&sentencesFromClient).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save sentences.",
		})
	}

	updatedSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph after updating.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": updatedSentences,
	})
}

func UpdateSentence(c *fiber.Ctx) error {
	var clientSentence models.Sentence
	paragraph := c.Query("paragraph")

	if len(paragraph) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No paragraph in query.",
		})
	}

	err := c.BodyParser(&clientSentence)

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

	var existingSentence models.Sentence

	err = database.DB.Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id").Where("sentence.id = ? AND template.user_id = ?", clientSentence.ID, userId).Find(&existingSentence).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Sentence does not exist in database.",
		})
	}

	err = database.DB.Save(&clientSentence).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save sentences.",
		})
	}

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

	err = database.DB.Delete(&models.Sentence{}, ids).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete sentences.",
		})
	}

	newSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

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

	err = actions.AddBulkSentences(sentencesFromClient, existingSentences)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to save bulk sentences.",
		})
	}

	// Re-assign sentences to what's now on the database.
	updatedSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}
