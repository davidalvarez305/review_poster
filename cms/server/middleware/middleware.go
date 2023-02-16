package middleware

import (
	"strings"

	"github.com/davidalvarez305/content_go/server/actions"
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

	return c.Next()
}
