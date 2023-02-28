package handlers

import (
	"fmt"

	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/davidalvarez305/content_go/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetWords(c *fiber.Ctx) error {
	words := &actions.Words{}

	sessionUserId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "User ID not found in session storage.",
		})
	}

	userId := c.Params("userId")

	if sessionUserId != userId {
		return c.Status(429).JSON(fiber.Map{
			"data": "Not allowed to access these resources.",
		})
	}

	err = words.GetWords(sessionUserId)

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

	if wordName == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = word.GetWordByName(wordName, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": word,
	})
}

func CreateWord(c *fiber.Ctx) error {
	var body actions.CreateWordInput
	word := &actions.Word{}
	synonyms := actions.Synonyms{}

	err := c.BodyParser(&body)

	if err != nil || len(body.Word) == 0 || len(body.Tag) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	word.Word = &models.Word{
		ID:     body.ID,
		Name:   body.Word,
		Tag:    body.Tag,
		UserID: body.UserID,
	}
	err = word.CreateWord()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	if len(body.Synonyms) > 0 {
		for i := 0; i < len(body.Synonyms); i++ {
			synonym := &models.Synonym{
				Synonym: body.Synonyms[i],
				WordID:  word.ID,
				Word:    word.Word,
			}
			synonyms = append(synonyms, synonym)
		}

		err = synonyms.CreateSynonyms()

		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"data": err.Error(),
			})
		}
	}
	words := &actions.Words{}
	err = words.GetWords(fmt.Sprintf("%+v", body.UserID))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": words,
	})
}

func UpdateWord(c *fiber.Ctx) error {
	word := &actions.Word{}

	err := c.BodyParser(&word)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = word.UpdateWord(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": word,
	})
}

func DeleteWord(c *fiber.Ctx) error {
	word_id := c.Query("word")
	word := &actions.Word{}

	if word_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = word.DeleteWord(userId, word_id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
