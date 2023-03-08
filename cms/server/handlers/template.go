package handlers

import (
	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetTemplates(c *fiber.Ctx) error {
	templates := &actions.Templates{}
	userId := c.Params("userId")

	err := templates.GetTemplates(userId)

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
	template := &actions.Template{}
	templates := &actions.Templates{}
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = template.CreateTemplate(userId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = templates.GetTemplates(userId)

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
	template := &actions.Template{}
	userId := c.Params("userId")

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = template.UpdateTemplate(userId)

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
	template := &actions.Template{}

	if templateId == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in params.",
		})
	}

	err := template.DeleteTemplate(userId, templateId)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
