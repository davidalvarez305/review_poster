package handlers

import (
	"errors"
	"os"
	"time"

	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/sessions"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var user models.User
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	err = actions.CreateUser(user)

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
	user, err := actions.GetUserFromSession(c)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	if user.Email == "" {
		return c.Status(401).JSON(fiber.Map{
			"data": errors.New("no user found"),
		})
	}

	user.AuthHeaderString = os.Getenv("AUTH_HEADER_STRING")

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func Logout(c *fiber.Ctx) error {
	user, err := actions.GetUserFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	if user.Email == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "User not found.",
		})
	}

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return err
	}

	err = sess.Destroy()

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
	var user models.User
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = actions.Login(user, c)

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
	var body models.User

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	user, err := actions.UpdateUser(body)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to update user.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	user, err := actions.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get user from session.",
		})
	}

	err = database.DB.Delete(&user).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete user.",
		})
	}

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get current session.",
		})
	}

	err = sess.Destroy()

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to logout.",
		})
	}

	return c.Status(204).JSON(fiber.Map{
		"data": "Deleted!",
	})
}

func ChangePassword(c *fiber.Ctx) error {

	// Handle Client Input
	type changePasswordInput struct {
		NewPassword string `json:"newPassword"`
	}
	code := c.Params("code")

	if len(code) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "No code sent in request.",
		})
	}

	// Initialize Structs
	var body changePasswordInput

	err := c.BodyParser(&body)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to parse body.",
		})
	}

	// Get User From Session
	user, err := actions.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get user from session.",
		})
	}

	// Retrieve Token from DB
	token, err := actions.GetToken(code)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Failed to get token.",
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
	updatedUser, err := actions.ChangePassword(user, body.NewPassword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to get user.",
		})
	}

	// Delete Token
	err = actions.DeleteToken(token)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": "Failed to delete token.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": updatedUser,
	})
}

func RequestChangePasswordCode(c *fiber.Ctx) error {
	user, err := actions.GetUserFromSession(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	err = actions.RequestChangePasswordCode(user)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": user,
	})
}
