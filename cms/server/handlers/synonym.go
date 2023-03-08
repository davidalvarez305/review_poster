package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateSynonym(c *fiber.Ctx) error {
	synonym := &actions.Synonym{}

	err := c.BodyParser(&synonym)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = synonym.CreateSynonym()

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": synonym,
	})
}

func UpdateSynonyms(c *fiber.Ctx) error {
	synonyms := &actions.Synonyms{}
	word := c.Query("word")

	if word == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := c.BodyParser(&synonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = synonyms.UpdateSynonyms(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func GetSelectedSynonyms(c *fiber.Ctx) error {
	synonyms := &actions.Synonyms{}
	word := c.Query("word")
	userId := c.Params("userId")

	if word == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := synonyms.GetSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func DeleteSynonym(c *fiber.Ctx) error {
	synonyms := &actions.Synonyms{}
	s := c.Query("synonyms")
	word := c.Query("word")
	userId := c.Params("userId")

	ids, err := utils.GetIds(s)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = synonyms.DeleteSynonyms(ids, word, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func BulkSynonymsPost(c *fiber.Ctx) error {
	synonyms := &actions.Synonyms{}
	word := c.Query("word")
	userId := c.Params("userId")

	if word == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := c.BodyParser(&synonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	existingSynonyms := &actions.Synonyms{}
	err = existingSynonyms.GetSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.

	var synonymsToDelete actions.Synonyms
	var synonymsToAdd actions.Synonyms

	err = synonymsToDelete.DeleteBulkSynonyms(synonyms, existingSynonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = synonymsToAdd.AddBulkSynonyms(synonyms, existingSynonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Re-assign synonyms to what's now on the database.
	err = synonyms.GetSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}
