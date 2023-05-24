package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetUserWords(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var words []models.Word

	err := database.DB.Where("user_id = ?", userId).Preload("Synonyms").Preload("User").Find(&words).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": words,
	})
}

func GetUserWord(c *fiber.Ctx) error {
	wordName := c.Params("wordName")
	userId := c.Params("userId")

	if len(wordName) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	var word models.Word

	err := database.DB.Preload("Synonyms").Where("name = ? AND user_id = ?", wordName, userId).First(&word).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch word by name.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": word,
	})
}

func CreateUserWord(c *fiber.Ctx) error {
	var input models.Word
	userId := c.Params("userId")

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	user, err := actions.GetUserFromSession(c)

	word := models.Word{
		ID:     input.ID,
		Name:   input.Name,
		Tag:    input.Tag,
		UserID: user.ID,
	}

	err = database.DB.Where("user_id = ?", user.ID).Save(&word).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create word and synonyms.",
		})
	}

	var updatedWords models.Word

	err = database.DB.Where("user_id = ? AND name = ?", userId, input.Name).First(&updatedWords).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words after saving.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedWords,
	})
}

func UpdateUserWord(c *fiber.Ctx) error {
	wordId := c.Params("id")
	userId := c.Params("userId")
	var word models.Word

	if len(wordId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	err := c.BodyParser(&word)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = database.DB.Where("user_id = ? AND id = ?", userId, wordId).Save(&word).First(&word).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to update word.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": word,
	})
}

func DeleteUserWord(c *fiber.Ctx) error {
	wordId := c.Params("id")
	userId := c.Params("userId")
	var word models.Word

	if len(wordId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	err := database.DB.Where("user_id = ? AND id = ?", userId, wordId).Delete(&word).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete word.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
