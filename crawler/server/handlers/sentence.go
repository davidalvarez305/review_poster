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
				"data": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": sentences,
		})
	}

	sentences, err := actions.GetAllSentences(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	err = actions.CreateSentences(sentences, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	paragraphs, err := actions.GetTemplatesAndParagraphs(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	updatedSentences, err := actions.UpdateSentences(sentencesFromClient, paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	sentences, err := actions.DeleteSentences(ids, paragraph, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
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
			"data": err.Error(),
		})
	}

	existingSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.

	err = actions.DeleteBulkSentences(sentencesFromClient, existingSentences)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = actions.AddBulkSentences(sentencesFromClient, existingSentences, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Re-assign sentences to what's now on the database.
	updatedSentences, err := actions.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}
