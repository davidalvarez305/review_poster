package handlers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CreatePosts(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword     string `json:"keyword"`
		ParentGroup string `json:"parent_group"`
	}
	products := &actions.AmazonSearchResultsPages{}

	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	err = products.CreateReviewPosts(body.Keyword, body.ParentGroup)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to create review posts.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
