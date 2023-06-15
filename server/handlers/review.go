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
	userId := c.Params("userId")

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

	paragraphs, err := actions.GetUserJoinedSentencesByParagraph(body.Template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences.",
		})
	}

	products, err := actions.CreateReviewPosts(body.Keyword, body.GroupName, words, paragraphs)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to create review posts.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": len(products),
	})
}

func TestCategories(c *fiber.Ctx) error {
	var categories []models.Category
	err := database.DB.Preload("SubCategories.ReviewPosts").Find(&categories).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to create review posts.",
		})
	}

	var deleteCategories []models.Category
	var deleteSubCategories []models.SubCategory

	for _, category := range categories {
		if len(category.SubCategories) == 0 {
			deleteCategories = append(deleteCategories, category)
			continue
		}

		for _, subCategory := range category.SubCategories {
			if len(subCategory.ReviewPosts) == 0 {
				deleteSubCategories = append(deleteSubCategories, *subCategory)
			}
		}
	}

	if len(deleteCategories) > 0 {
		err = database.DB.Delete(&deleteCategories).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to create review posts.",
			})
		}
	}

	if len(deleteSubCategories) > 0 {
		err = database.DB.Delete(&deleteSubCategories).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to create review posts.",
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"data": "EZ.",
	})
}
