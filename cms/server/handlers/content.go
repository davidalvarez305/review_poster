package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	template := c.Query("template")
	userId := c.Params("userId")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No query string found in the URL.",
		})
	}

	sentences, err := actions.GetSentences(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to query sentences.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func GetDictionary(c *fiber.Ctx) error {
	userId := c.Params("userId")
	dict, err := actions.GetDictionary(userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to query dictionary.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dict,
	})
}
