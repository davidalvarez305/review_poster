package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CrawlAmazon(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	data, err := actions.ScrapeSearchResultsPage(body.Keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func SearchPAAPI5(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	p, err := actions.SearchPaapi5Items(body.Keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": p,
	})
}
