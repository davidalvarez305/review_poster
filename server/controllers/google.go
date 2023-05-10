package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Google(router fiber.Router) {
	google := router.Group("google")

	google.Get("/keywords", handlers.GetCommercialKeywords)
	google.Get("/seed", handlers.GetSeedKeywords)
}
