package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/gofiber/fiber/v2"
)

func ParseTags(c *fiber.Ctx) error {
	tag := c.Query("tag")

	if len(tag) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Tag missing in request.",
		})
	}

	tagRegex := regexp.MustCompile(`^\(#\w+\)$`)

	if len(tagRegex.FindAllString(tag, 3)) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Incorrectly formatted tag.",
		})
	}

	prompt := fmt.Sprintf("Consider the following tags:\n%s.\nTogether they make a full sentence. Can you parse each tag, and list 50 different ways that each one of these tags can be written, in a way that makes them grammatically correct sentences? Please don't enumerate the list, return it as a string separated by lines.", tag)

	response, err := actions.QueryOpenAI(prompt)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to query Open AI.",
		})
	}

	if len(response.Choices) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No data return from Open AI.",
		})
	}

	var parsedPhrases []string

	for _, phrase := range strings.Split(response.Choices[0].Text, "\n") {
		if len(phrase) > 0 {
			parsedPhrases = append(parsedPhrases, phrase)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"data": parsedPhrases,
	})
}
