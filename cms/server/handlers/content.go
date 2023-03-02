package handlers

import (
	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	sentences := &actions.Sentences{}
	template := c.Query("template")
	userId := c.Params("userId")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No query string found in the URL.",
		})
	}

	err := sentences.GetSentences(template, userId)

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
	words := &actions.Words{}
	userId := c.Params("userId")
	err := words.GetWords(userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dictionary,
	})
}
