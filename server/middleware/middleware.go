package middleware

import (
	"os"
	"strings"

	"github.com/davidalvarez305/review_poster/server/actions"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	auth := headers["Authorization"]
	path := c.Path()

	if path == "/api/user/login" || path == "/api/user/register" {
		return c.Next()
	}

	if auth != "Bearer "+os.Getenv("AUTH_HEADER_STRING") {
		return c.Status(401).JSON(fiber.Map{
			"data": "Unauthorized request.",
		})
	}

	return c.Next()
}

func ResourceAccessRestriction(c *fiber.Ctx) error {
	protectedRoutes := []string{"word", "template", "sentence", "paragraph", "synonym", "content", "dictionary"}
	var found []string
	path := c.OriginalURL()
	headers := c.GetReqHeaders()
	secretAgent := headers["X-Secret-Agent"]

	for _, route := range protectedRoutes {
		if strings.Contains(path, route) {
			found = append(found, route)
		}
	}

	// If request was not made to one of the protected routes, continue on.
	// This is for login, register, etc...
	if len(found) == 0 {
		return c.Next()
	}

	if os.Getenv("PRODUCTION") == "0" || secretAgent == os.Getenv("X_SECRET_AGENT") {
		return c.Next()
	}

	sessionUserId, err := actions.GetUserIdFromSession(c)

	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"data": "User ID not found in session storage.",
		})
	}

	userId := c.Params("userId")

	if sessionUserId != userId {
		return c.Status(403).JSON(fiber.Map{
			"data": "Not allowed to access these resources.",
		})
	}

	return c.Next()
}
