package handlers

import (
	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	content := &actions.Content{}
	template := c.Query("template")

	if template == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No query string found in the URL.",
		})
	}

	err := content.GetContent(template)

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
	dictionary := &actions.Dictionary{}
	err := dictionary.GetDictionary()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": dictionary,
	})
}
