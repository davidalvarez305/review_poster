package handlers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CrawlAmazon(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	if len(keyword) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No keyword found in query string.",
		})
	}

	products, err := actions.ScrapeSearchResultsPage(keyword)

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
	keyword := c.Query("keyword")

	if len(keyword) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No keyword found in query string.",
		})
	}

	products, err := actions.SearchPaapi5Items(keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to search PAAPI5.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
