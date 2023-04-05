package handlers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/davidalvarez305/review_poster/crawler/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	sentences, err := actions.PullDynamicContent()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func GetDictionary(c *fiber.Ctx) error {
	data, err := actions.PullContentDictionary()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func GetDynamicContent(c *fiber.Ctx) error {
	productName := c.Query("productName")

	if len(productName) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No product name in query string.",
		})
	}

	data, err := actions.PullContentDictionary()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	sentences, err := actions.PullDynamicContent()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	content := utils.GenerateContentUtil(productName, data.Data, sentences.Data)

	return c.Status(200).JSON(fiber.Map{
		"data": content,
	})
}

func GetOpenAIContent(c *fiber.Ctx) error {
	prompt := c.Query("prompt")

	if len(prompt) == 0 {
		return c.Status(200).JSON(fiber.Map{
			"data": "No query in URL.",
		})
	}

	data, err := actions.GetAdditionalContent(prompt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}
