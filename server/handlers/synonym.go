package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateSynonym(c *fiber.Ctx) error {
	var synonym models.Synonym

	err := c.BodyParser(&synonym)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&synonym).First(&synonym).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save synonym.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": synonym,
	})
}

func UpdateSynonyms(c *fiber.Ctx) error {
	var synonyms []models.Synonym
	word := c.Query("word")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := c.BodyParser(&synonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user ID from session.",
		})
	}

	err = database.DB.Save(&synonyms).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save synonyms.",
		})
	}

	updatedSynonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSynonyms,
	})
}

func UpdateSynonym(c *fiber.Ctx) error {
	var clientSynonym models.Synonym
	word := c.Query("word")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := c.BodyParser(&clientSynonym)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user ID from session.",
		})
	}

	var existingSynonym models.Synonym

	err = database.DB.Joins("JOIN word ON word.id = synonym.word_id").Where("synonym.id = ? AND word.user_id = ?", clientSynonym.ID, userId).Find(&existingSynonym).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Synonym does not exist in database.",
		})
	}

	err = database.DB.Save(&clientSynonym).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save synonyms.",
		})
	}

	updatedSynonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSynonyms,
	})
}

func GetSelectedSynonyms(c *fiber.Ctx) error {
	word := c.Query("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	synonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func GetUserSynonymsByWord(c *fiber.Ctx) error {
	word := c.Params("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	synonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user's synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func DeleteSynonym(c *fiber.Ctx) error {
	s := c.Query("synonyms")
	word := c.Query("word")
	userId := c.Params("userId")

	ids, err := utils.GetIds(s)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to parse IDs.",
		})
	}

	err = database.DB.Delete(&models.Synonym{}, ids).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete synonyms.",
		})
	}

	synonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": synonyms,
	})
}

func BulkSynonymsPost(c *fiber.Ctx) error {
	var clientSynonyms []models.Synonym
	word := c.Query("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in query.",
		})
	}

	err := c.BodyParser(&clientSynonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	existingSynonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	// These functions will filter synonyms coming from the client & compare with existing ones.
	// It will keep anything that's new, and delete what was not sent from the client.
	err = actions.DeleteBulkSynonyms(clientSynonyms, existingSynonyms)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete bulk synonyms.",
		})
	}

	err = actions.AddBulkSynonyms(clientSynonyms, existingSynonyms)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to add bulk synonyms.",
		})
	}

	// Re-assign synonyms to what's now on the database.
	updatedSynonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSynonyms,
	})
}
