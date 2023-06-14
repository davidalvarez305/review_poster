package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func Google(router fiber.Router) {
	google := router.Group("google", middleware.AuthMiddleware)

	google.Get("/keywords", handlers.GetCommercialKeywords)
	google.Get("/seed", handlers.GetSeedKeywords)
	google.Get("/generate", handlers.GenerateKeywords)
}
