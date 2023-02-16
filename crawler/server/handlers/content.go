package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/utils"
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
	type reqBody struct {
		ProductName string `json:"productName"`
	}
	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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

	content := utils.GenerateContentUtil(body.ProductName, data, sentences)

	return c.Status(200).JSON(fiber.Map{
		"data": content,
	})
}
