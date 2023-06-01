package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/utils"
	"github.com/gofiber/fiber/v2"
)

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
	sentences, err := actions.GetSentencesByTemplate(template, userId)

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
