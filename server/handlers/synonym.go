package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
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

	if len(input.DeleteSynonyms) > 0 {
		err = database.DB.Delete(&models.Synonym{}, input.DeleteSynonyms).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to delete synonyms.",
			})
		}
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
