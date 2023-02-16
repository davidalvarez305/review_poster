package handlers

import (
	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/davidalvarez305/content_go/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetSentencesByParagraph(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	paragraph := c.Params("paragraph")

	if paragraph == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No paragraph in params.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

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

func GetSentences(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentences.GetAllSentences(userId)

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
	jp := &actions.JoinedParagraphs{}

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = sentences.CreateSentences()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = jp.GetTemplatesAndParagraphs(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": jp,
	})
}

func UpdateSentences(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	paragraph := c.Query("paragraph")

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

	userId, err := actions.GetUserIdFromSession(c)

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

	if paragraph == "" || s == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No paragraph or sentences in query.",
		})
	}

	ids, err := utils.GetIds(s)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
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

func GetTemplatesAndParagraphs(c *fiber.Ctx) error {
	jp := &actions.JoinedParagraphs{}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = jp.GetTemplatesAndParagraphs(userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": jp,
	})
}

func BulkSentencesUpdate(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	paragraph := c.Query("paragraph")

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

	userId, err := actions.GetUserIdFromSession(c)

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

	err = sentencesToAdd.AddBulkSentences(sentences, existingSentences)

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
