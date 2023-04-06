package routes

import (
	"github.com/davidalvarez305/review_poster/cms/server/controllers"
	"github.com/davidalvarez305/review_poster/cms/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("api", middleware.AuthMiddleware, middleware.ResourceAccessRestriction)
	controllers.User(api)
}
