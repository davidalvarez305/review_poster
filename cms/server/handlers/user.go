package handlers

import (
	"errors"
	"time"

	"github.com/davidalvarez305/content_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	user := &actions.Users{}
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = user.CreateUser()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": user,
	})
}

func GetUser(c *fiber.Ctx) error {
	user := &actions.Users{}

	err := user.GetUserFromSession(c)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	if user.Users == nil {
		return c.Status(404).JSON(fiber.Map{
			"data": errors.New("no user found"),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func Logout(c *fiber.Ctx) error {
	user := &actions.Users{}

	err := user.Logout(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": "Logged out!",
	})
}

func Login(c *fiber.Ctx) error {
	user := &actions.Users{}
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = user.Login(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	var body actions.Users
	user := &actions.Users{}

	err := c.BodyParser(&body)

	if err != nil {
		return err
	}

	err = user.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = user.UpdateUser(body)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	user := &actions.Users{}

	err := user.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = user.Delete()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = user.Logout(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}

func ChangePassword(c *fiber.Ctx) error {

	// Handle Client Input
	type ChangePasswordInput struct {
		NewPassword string `json:"newPassword"`
	}
	code := c.Params("code")

	if code == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "No code sent in request.",
		})
	}

	// Initialize Structs
	var body ChangePasswordInput
	user := &actions.Users{}
	token := &actions.Token{}

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Get User From Session
	err = user.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Retrieve Token from DB
	err = token.GetToken(code, user.ID)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Validate Token's Expiry Date
	difference := time.Now().Unix() - token.CreatedAt

	if difference > 300 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Token expired.",
		})
	}

	// Update User
	err = user.ChangePassword(body.NewPassword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	// Delete Token
	err = token.DeleteToken()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func RequestChangePasswordCode(c *fiber.Ctx) error {
	user := &actions.Users{}

	err := user.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = user.RequestChangePasswordCode()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}
