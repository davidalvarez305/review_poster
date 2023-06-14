package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func DynamicContent(router fiber.Router) {

	content := router.Group("content", middleware.AuthMiddleware)
	content.Get("/dynamic", handlers.GetDynamicContent)
}
