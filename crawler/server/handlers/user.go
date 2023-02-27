package handlers

import (
	"fmt"
	"os"

	"github.com/davidalvarez305/review_poster/crawler/server/actions"
	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
	"github.com/davidalvarez305/review_poster/crawler/server/sessions"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	var u types.User
	var user types.UserRequestBody
	err := c.BodyParser(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Unable to Parse Request Body.",
		})
	}

	if user.Password == "" || user.Username == "" {
		return c.Status(400).JSON(fiber.Map{
			"data": "Missing Fields.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	u.Username = user.Username
	u.Password = hashedPassword

	fmt.Printf("User %v", u)

	data, err := actions.CreateUser(u)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"data": data,
	})
}

func GetUser(c *fiber.Ctx) error {
	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	k := sess.Get(os.Getenv("COOKIE_NAME"))

	if k == nil {
		return c.Status(404).JSON(fiber.Map{
			"data": "Not found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": k,
	})
}

func Logout(c *fiber.Ctx) error {

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	k := sess.Get(os.Getenv("COOKIE_NAME"))

	if k == nil {
		return c.Status(404).JSON(fiber.Map{
			"data": "Not found.",
		})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": "Logged out!",
	})
}

func Login(c *fiber.Ctx) error {
	var u types.UserRequestBody
	var user models.User
	err := c.BodyParser(&u)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	result := database.DB.Where("username = ?", &u.Username).First(&user)

	if result.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"data": "Incorrect username.",
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "Incorrect password.",
		})
	}

	id := sessions.Sessions.KeyGenerator()

	sess, err := sessions.Sessions.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	sess.Set(os.Getenv("COOKIE_NAME"), id)

	if err := sess.Save(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": id,
	})
}
