package handlers

import (
	"strings"

	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/gofiber/fiber/v2"
)

func GetCommercialKeywords(c *fiber.Ctx) error {
	body := &types.Keyword{}
	results := &actions.GoogleKeywordResults{}
	err := c.BodyParser(body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	s := strings.Split(body.Keyword, "\n")

	if len(s) > 1 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Only one seed keyword allowed per query.",
		})
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{body.Keyword},
		},
	}

	err = results.QueryGoogle(q)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to query Google.",
		})
	}

	if len(results.Results) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	seedKeywords, err := actions.GetSeedKeywords(results)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get seed keywords.",
		})
	}

	if len(seedKeywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Seed Keywords Found.",
		})
	}

	keywords, err := actions.GetCommercialKeywords(seedKeywords)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get commercial keywords.",
		})
	}

	if len(keywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Commercial Keywords Found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": keywords,
	})
}

func GetSeedKeywords(c *fiber.Ctx) error {
	body := &types.Keyword{}
	results := &actions.GoogleKeywordResults{}
	err := c.BodyParser(body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse request body.",
		})
	}

	s := strings.Split(body.Keyword, "\n")

	if len(s) > 1 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Only one seed keyword allowed per query.",
		})
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{body.Keyword},
		},
	}

	err = results.QueryGoogle(q)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to query Google.",
		})
	}

	if len(results.Results) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	seedKeywords, err := actions.GetSeedKeywords(results)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get seed keywords.",
		})
	}

	if len(seedKeywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Seed Keywords Found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": seedKeywords,
	})
}
