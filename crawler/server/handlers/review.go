package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CreatePosts(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword     string `json:"keyword"`
		ParentGroup string `json:"parent_group"`
	}

	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	products, err := actions.CreateReviewPosts(body.Keyword, body.ParentGroup)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
