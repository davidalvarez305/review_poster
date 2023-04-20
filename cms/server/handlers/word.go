package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetWords(c *fiber.Ctx) error {
	wordName := c.Query("word")
	userId := c.Params("userId")

	if len(wordName) > 0 {
		word, err := actions.GetWordByName(wordName, userId)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": "Failed to query word by name.",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": word,
		})
	}

	words, err := actions.GetWords(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": words,
	})
}

func CreateWord(c *fiber.Ctx) error {
	var body actions.CreateWordInput
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

	word := models.Word{
		ID:     body.ID,
		Name:   body.Word,
		Tag:    body.Tag,
		UserID: body.UserID,
	}

	for i := 0; i < len(body.Synonyms); i++ {
		synonym := &models.Synonym{
			Synonym: body.Synonyms[i],
			WordID:  word.ID,
			Word:    &word,
		}
		word.Synonyms = append(word.Synonyms, synonym)
	}

	err = actions.CreateWord(word)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create word and synonyms.",
		})
	}

	words, err := actions.GetWords(userId)

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
	var word models.Word

	userId := c.Params("userId")

	err := c.BodyParser(&word)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = actions.UpdateWord(word, userId)

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
	var word models.Word

	if wordId == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := actions.DeleteWord(word, userId, wordId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete word.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
