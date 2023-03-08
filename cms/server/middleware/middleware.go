package middleware

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/review_poster/cms/server/actions"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	url := c.OriginalURL()

	if strings.Contains(url, "/login") || strings.Contains(url, "/register") {
		return c.Next()
	}

	userId, err := actions.GetUserIdFromSession(c)

	if err != nil || len(userId) == 0 {
		return c.Status(401).JSON(fiber.Map{
			"data": "Unauthorized.",
		})
	}

	fmt.Printf("%+v", c.Request())

	return c.Next()
}

func ResourceAccessRestriction(c *fiber.Ctx) error {
	sessionUserId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"data": "User ID not found in session storage.",
		})
	}

	userId := c.Params("userId")

	if sessionUserId != userId {
		return c.Status(429).JSON(fiber.Map{
			"data": "Not allowed to access these resources.",
		})
	}

	return c.Next()
}
