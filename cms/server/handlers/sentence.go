package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetSentences(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	userId := c.Params("userId")
	paragraph := c.Query("paragraph")

	if len(paragraph) > 0 {
		err := sentences.GetSentencesByParagraph(paragraph, userId)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": sentences,
		})
	}

	err := sentences.GetAllSentences(userId)

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
	sentences := &actions.Sentences{}
	paragraphs := &actions.Paragraphs{}
	userId := c.Params("userId")

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentences.CreateSentences(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = paragraphs.GetTemplatesAndParagraphs(userId)

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
	sentences := &actions.Sentences{}
	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentences.UpdateSentences(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": sentences,
	})
}

func DeleteSentence(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	s := c.Query("sentences")
	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	ids, err := utils.GetIds(s)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentences.DeleteSentences(ids, paragraph, userId)

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
	sentences := &actions.Sentences{}
	paragraph := c.Query("paragraph")
	userId := c.Params("userId")

	if paragraph == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No paragraph in query.",
		})
	}

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	existingSentences := &actions.Sentences{}
	err = existingSentences.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.

	var sentencesToDelete actions.Sentences
	var sentencesToAdd actions.Sentences

	err = sentencesToDelete.DeleteBulkSentences(sentences, existingSentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentencesToAdd.AddBulkSentences(sentences, existingSentences, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Re-assign sentences to what's now on the database.
	err = sentences.GetSentencesByParagraph(paragraph, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}
