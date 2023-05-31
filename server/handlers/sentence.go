package handlers

import (
	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
	"github.com/gofiber/fiber/v2"
)

func GetUserParagraphSentencesByTemplate(c *fiber.Ctx) error {
	userId := c.Params("userId")
	paragraph := c.Params("paragraphName")
	template := c.Params("templateName")

	if len(template) == 0 || len(paragraph) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing URL params.",
		})
	}

	sentences, err := actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func CreateUserParagraphSentencesByTemplate(c *fiber.Ctx) error {
	var sentences []models.Sentence

	userId := c.Params("userId")
	paragraph := c.Params("paragraphName")
	template := c.Params("templateName")

	if len(template) == 0 || len(paragraph) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing URL params.",
		})
	}

	err := c.BodyParser(&sentences)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&sentences).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to create sentences.",
		})
	}

	sentences, err = actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": sentences,
	})
}

func UpdateUserParagraphSentencesByTemplate(c *fiber.Ctx) error {
	paragraph := c.Params("paragraphName")
	template := c.Params("templateName")
	userId := c.Params("userId")

	if len(template) == 0 || len(paragraph) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing URL params.",
		})
	}

	var input types.UpdateUserParagraphSentencesByTemplateInput

	err := c.BodyParser(&input)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	err = database.DB.Save(&input.Sentences).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save sentences.",
		})
	}

	updatedSentences, err := actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch sentences by paragraph after updating.",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": updatedSentences,
	})
}

func DeleteUserParagraphSentencesByTemplate(c *fiber.Ctx) error {
	paragraph := c.Params("paragraphName")
	template := c.Params("templateName")
	userId := c.Params("userId")

	if len(template) == 0 || len(paragraph) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing URL params.",
		})
	}

	var input types.DeleteUserParagraphSentencesByTemplateInput

	if len(input.DeleteSentences) > 0 {
		err := database.DB.Delete(&models.Synonym{}, input.DeleteSentences).Error

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"data": "Failed to delete synonyms.",
			})
		}
	}

	updatedSentences, err := actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}

func UpdateUserSentence(c *fiber.Ctx) error {
	paragraph := c.Params("paragraphName")
	template := c.Params("templateName")
	userId := c.Params("userId")
	sentenceId := c.Params("sentenceId")

	if len(template) == 0 || len(paragraph) == 0 || len(sentenceId) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing URL params.",
		})
	}

	var sentence models.Sentence

	err := c.BodyParser(&sentence)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	userId, err = actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch user ID from session.",
		})
	}

	var existingSentence models.Sentence

	err = database.DB.Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id").Where("sentence.id = ? AND template.user_id = ? AND paragraph.name = ?", sentenceId, userId, paragraph).Find(&existingSentence).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Sentence does not exist in database.",
		})
	}

	existingSentence.Sentence = sentence.Sentence
	existingSentence.ParagraphID = sentence.ParagraphID

	err = database.DB.Save(&existingSentence).Error

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to save sentence.",
		})
	}

	updatedSentences, err := actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to fetch synonyms by word.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}

func DeleteUserSentence(c *fiber.Ctx) error {
	template := c.Params("template")
	userId := c.Params("userId")
	paragraph := c.Params("paragraphName")
	sentenceId := c.Params("sentenceId")

	if len(template) == 0 || len(paragraph) == 0 || len(sentenceId) == 0 {
		return c.Status(500).JSON(fiber.Map{
			"data": "Template or paragraph not in URL params.",
		})
	}

	var existingSentence models.Sentence

	err := database.DB.Joins("INNER JOIN paragraph ON paragraph.id = sentence.paragraph_id INNER JOIN template ON template.id = paragraph.template_id").Where("template.name = ? AND template.user_id = ? AND id = ?", template, userId, sentenceId).First(&existingSentence).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to find sentence.",
		})
	}

	err = database.DB.Delete(&existingSentence).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete sentence.",
		})
	}

	updatedSentences, err := actions.GetUserSentencesByParagraphAndTemplate(paragraph, userId, template)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to fetch sentences by template.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedSentences,
	})
}
