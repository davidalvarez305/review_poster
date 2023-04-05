package routes

import (
	"github.com/davidalvarez305/review_poster/crawler/server/controllers"
	"github.com/davidalvarez305/review_poster/crawler/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("api", middleware.AuthMiddleware)
	controllers.Google(api)
	controllers.Amazon(api)
	controllers.ReviewPost(api)
	controllers.DynamicContent(api)
}
