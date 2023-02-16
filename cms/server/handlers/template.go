package handlers

import (
	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func GetTemplates(c *fiber.Ctx) error {
	templates := &actions.Templates{}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Get Selected Synonyms.",
		})
	}

	err = templates.GetTemplates(userId)

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

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = template.CreateTemplate()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
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

	err := c.BodyParser(&template)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
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
	template_id := c.Query("template")
	template := &actions.Template{}

	if template_id == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No template in query.",
		})
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = template.DeleteTemplate(userId, template_id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}
