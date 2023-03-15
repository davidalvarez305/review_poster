package handlers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/gofiber/fiber/v2"
)

func CrawlAmazon(c *fiber.Ctx) error {
	body := &types.Keyword{}
	err := c.BodyParser(body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	products, err := actions.ScrapeSearchResultsPage(body.Keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to scrape Amazon SERP.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}

func SearchPAAPI5(c *fiber.Ctx) error {
	body := &types.Keyword{}
	err := c.BodyParser(body)
	products := &actions.PAAAPI5Response{}

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	err = products.SearchPaapi5Items(body.Keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to search PAAPI5.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
