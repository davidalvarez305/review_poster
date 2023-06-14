package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func Ai(router fiber.Router) {

	ai := router.Group("ai", middleware.AuthMiddleware)
	ai.Get("/tags", handlers.ParseTags)
}
