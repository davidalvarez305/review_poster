package routes

import (
	"github.com/davidalvarez305/review_poster/cms/server/controllers"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	// middleware.AuthMiddleware, middleware.ResourceAccessRestriction
	api := app.Group("api")
	controllers.User(api)
}
