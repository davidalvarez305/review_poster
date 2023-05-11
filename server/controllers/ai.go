package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Ai(router fiber.Router) {

	ai := router.Group("ai")
	ai.Get("/tags", handlers.ParseTags)
}
