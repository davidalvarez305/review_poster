package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
	"github.com/gofiber/fiber/v2"
)

func UpdateUserSynonymByWord(c *fiber.Ctx) error {
	var synonym models.Synonym
	word := c.Params("word")
	synonymId := c.Params("synonymId")

	if len(word) == 0 || len(synonymId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	err := c.BodyParser(&synonym)

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

	err = database.DB.Joins("JOIN word ON word.id = synonym.word_id").Where("synonym.id = ? AND word.user_id = ? AND word.name = ?", synonym.ID, userId, word).Find(&existingSynonym).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Synonym does not exist in database.",
		})
	}

	existingSynonym.Synonym = synonym.Synonym
	existingSynonym.WordID = synonym.WordID

	err = database.DB.Save(&existingSynonym).Error

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

func CreateUserSynonymsByWord(c *fiber.Ctx) error {
	word := c.Params("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	var synonyms []models.Synonym

	err := c.BodyParser(&synonyms)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&synonyms).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to save synonyms.",
		})
	}

	createdSynonyms, err := actions.GetUserSynonymsByWord(word, userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user's synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": createdSynonyms,
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

func DeleteUserSynonymByWord(c *fiber.Ctx) error {
	word := c.Params("word")
	userId := c.Params("userId")
	synonymId := c.Params("synonymId")

	if len(word) == 0 || len(synonymId) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"data": "Word or synonym not in URL params.",
		})
	}

	var synonym models.Synonym

	err := database.DB.Joins("INNER JOIN word ON word.id = synonym.word_id").Where("word.name = ? AND word.user_id = ? AND synonym.id = ?", word, userId, synonymId).First(&synonym).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to find synonym.",
		})
	}

	err = database.DB.Delete(&synonym).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete synonym.",
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

func UpdateUserSynonymsByWord(c *fiber.Ctx) error {
	word := c.Params("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	var input types.UpdateUserSynonymsByWordInput

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	if len(input.Synonyms) > 0 {
		err = database.DB.Save(&input.Synonyms).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to save synonyms.",
			})
		}
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

func DeleteUserSynonymsByWord(c *fiber.Ctx) error {
	word := c.Params("word")
	userId := c.Params("userId")

	if len(word) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No word in URL params.",
		})
	}

	var input types.DeleteUserSynonymsByWordInput

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	if len(input.DeleteSynonyms) > 0 {
		err = database.DB.Delete(&models.Synonym{}, input.DeleteSynonyms).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to delete synonyms.",
			})
		}
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
