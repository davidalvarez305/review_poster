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

	if os.Getenv("PRODUCTION") == "0" {
		return c.Next()
	}

	if auth != "Bearer "+os.Getenv("AUTH_HEADER_STRING") {
		return c.Status(401).JSON(fiber.Map{
			"data": "Unauthorized request.",
		})
	}

	return c.Next()
}

func PostMiddleware(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	auth := headers["Authorization"]
	method := c.Method()
	path := c.Path()

	if method == "GET" || path == "/api/user/login" || path == "/api/user/register" || os.Getenv("PRODUCTION") == "0" {
		return c.Next()
	}

	token, err := actions.GetCsrfTokenFromSession(c)

	if auth != token || err != nil {
		return c.Status(401).JSON(fiber.Map{
			"data": "Unauthorized request.",
		})
	}

	return c.Next()
}

func ResourceAccessRestriction(fn fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			return fn(c)
		}

		// This is for testing in dev
		if secretAgent == os.Getenv("X_SECRET_AGENT") {
			return fn(c)
		}

		sessionUser, err := actions.GetUserFromSession(c)

		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"data": "User ID not found in session storage.",
			})
		}

		userId, err := c.ParamsInt("userId")

		if err != nil {
			return c.Status(403).JSON(fiber.Map{
				"data": "Could not convert ID from params to int.",
			})
		}

		if sessionUser.ID != userId {
			return c.Status(403).JSON(fiber.Map{
				"data": "Not allowed to access these resources.",
			})
		}

		return fn(c)
	}
}
