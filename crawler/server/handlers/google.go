package handlers

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/gofiber/fiber/v2"
)

func GetCommercialKeywords(c *fiber.Ctx) error {
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

	results, err := actions.QueryGoogle(q)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
		})
	}

	fmt.Printf("Length of Seed Keywords: %v", seedKeywords)

	if len(seedKeywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Seed Keywords Found.",
		})
	}

	keywords, err := actions.GetCommercialKeywords(seedKeywords)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	fmt.Printf("Length of Commercial Keywords: %v", keywords)

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
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
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

	results, err := actions.QueryGoogle(q)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
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
			"data": err.Error(),
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

func GetPeopleAlsoAsk(c *fiber.Ctx) error {

	keywords, err := actions.CrawlGoogleSERP("bbcor bat reviews")

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": keywords,
	})
}
