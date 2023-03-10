package controllers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Google(router fiber.Router) {
	google := router.Group("google")

	google.Post("/keywords", handlers.GetCommercialKeywords)
	google.Post("/seed", handlers.GetSeedKeywords)
}
