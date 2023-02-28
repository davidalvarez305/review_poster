package handlers

import (
	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/davidalvarez305/content_go/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetWords(c *fiber.Ctx) error {
	words := &actions.Words{}

	userId := c.Params("userId")

	err := words.GetWords(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Faield to query words.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": words,
	})
}

func GetWord(c *fiber.Ctx) error {
	word := &actions.Word{}
	wordName := c.Params("word")
	userId := c.Params("userId")

	err := word.GetWordByName(wordName, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query word by name.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": word,
	})
}

func CreateWord(c *fiber.Ctx) error {
	var body actions.CreateWordInput
	word := &actions.Word{}
	userId := c.Params("userId")

	err := c.BodyParser(&body)

	if err != nil || len(body.Word) == 0 || len(body.Tag) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	if len(body.Synonyms) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Must have synonyms to create word.",
		})
	}

	word.Word = &models.Word{
		ID:     body.ID,
		Name:   body.Word,
		Tag:    body.Tag,
		UserID: body.UserID,
	}

	for i := 0; i < len(body.Synonyms); i++ {
		synonym := &models.Synonym{
			Synonym: body.Synonyms[i],
			WordID:  word.ID,
			Word:    word.Word,
		}
		word.Synonyms = append(word.Synonyms, synonym)
	}

	err = word.CreateWord()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create word and synonyms.",
		})
	}

	words := &actions.Words{}

	err = words.GetWords(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words after saving.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": words,
	})
}

func UpdateWord(c *fiber.Ctx) error {
	word := &actions.Word{}

	userId := c.Params("userId")

	err := c.BodyParser(&word)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = word.UpdateWord(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to update word.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": word,
	})
}

func DeleteWord(c *fiber.Ctx) error {
	wordId := c.Params("id")
	userId := c.Params("userId")
	word := &actions.Word{}

	if wordId == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := word.DeleteWord(userId, wordId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete word.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
