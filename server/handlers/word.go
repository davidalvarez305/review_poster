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
		var word models.Word

		err := database.DB.Where("name = ? AND user_id = ?", wordName, userId).Find(&word).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to fetch word by name.",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": word,
		})
	}

	var words []models.Word

	err := database.DB.Where("user_id = ?", userId).Preload("User").Find(&words).Error

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
	
	get user from session, wtf

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

	err = database.DB.Where("user_id = ?", word.UserID).Save(&word).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create word and synonyms.",
		})
	}

	var updatedWords []models.Word

	err = database.DB.Where("user_id = ?", userId).Find(&updatedWords).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words after saving.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedWords,
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

	err = database.DB.Where("user_id = ? AND id = ?", userId, word.ID).Save(&word).First(&word).Error

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

	if len(wordId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in querystring.",
		})
	}

	err := database.DB.Where("user_id = ? AND id = ?", userId, word_id).Delete(&word).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete word.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
