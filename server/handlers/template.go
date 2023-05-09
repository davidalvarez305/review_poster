package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/gofiber/fiber/v2"
)

func GetTemplates(c *fiber.Ctx) error {
	userId := c.Params("userId")

	templates, err := actions.GetTemplates(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": templates,
	})
}

func CreateTemplate(c *fiber.Ctx) error {
	var template models.Template
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = actions.CreateTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	templates, err := actions.GetTemplates(userId)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": templates,
	})
}

func UpdateTemplate(c *fiber.Ctx) error {
	var template models.Template
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = actions.UpdateTemplate(template, userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": template,
	})
}

func DeleteTemplate(c *fiber.Ctx) error {
	templateId := c.Params("templateId")
	userId := c.Params("userId")
	var template models.Template

	if templateId == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in params.",
		})
	}

	err := actions.DeleteTemplate(template, userId, templateId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
