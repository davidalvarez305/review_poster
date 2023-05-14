package handlers

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	template := c.Query("template")
	userId := c.Params("userId")

	if len(template) == 0 || len(userId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Incorrect query or URL params.",
		})
	}

	var sentences []models.Sentence

	err := database.DB.Where("\"Template\".user_id = ? AND \"Template\".name = ?", userId, template).Joins("Template").Preload("Paragraph").Find(&sentences).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func GetDynamicContent(c *fiber.Ctx) error {
	productName := c.Query("productName")
	template := c.Query("template")
	userId := c.Params("userId")

	if len(template) == 0 || len(productName) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Incorrect querystring.",
		})
	}

	// First get the words (dictionary)
	var words []models.Word

	err := database.DB.Where("user_id = ?", userId).Preload("Synonyms").Find(&words).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch words.",
		})
	}

	// Second get the sentences
	var sentences []models.Sentence

	err = database.DB.Where("\"Template\".user_id = ? AND \"Template\".name = ?", userId, template).Joins("Template").Preload("Paragraph").Find(&sentences).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences.",
		})
	}

	content := utils.GenerateContentUtil(productName, words, sentences)

	return c.Status(200).JSON(fiber.Map{
		"data": content,
	})
}
