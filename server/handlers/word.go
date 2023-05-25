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

	var createdWord models.Word

	err = database.DB.Where("user_id = ? AND name = ?", userId, input.Name).First(&createdWord).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query words after saving.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": createdWord,
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

	err = database.DB.Where("user_id = ?", userId).Save(&word).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save word.",
		})
	}

	err = database.DB.Where("user_id = ? AND id = ?", userId, wordId).First(&word).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch updated word.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": word,
	})
}
