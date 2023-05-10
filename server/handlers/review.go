package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
	"github.com/gofiber/fiber/v2"
)

func CreatePosts(c *fiber.Ctx) error {
	var body types.CreateReviewPostsInput

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	if len(body.GroupName) == 0 || len(body.Keyword) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Either group or keyword is missing in request body.",
		})
	}

	var words []models.Word

	err = database.DB.Preload("Synonyms").Find(&words).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch dictionary data.",
		})
	}

	sentences, err := actions.GetSentencesByTemplate(body.Template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences.",
		})
	}

	products, err := actions.CreateReviewPosts(body.Keyword, body.GroupName, words, sentences)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to create review posts.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": len(products),
	})
}
