package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	content := &actions.Content{}
	template := c.Query("template")
	userId := c.Params("userId")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No query string found in the URL.",
		})
	}

	err := content.GetSentences(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": content,
	})
}

func GetDictionary(c *fiber.Ctx) error {
	dict := &actions.Dictionary{}
	userId := c.Params("userId")
	err := dict.GetDictionary(userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dict,
	})
}
