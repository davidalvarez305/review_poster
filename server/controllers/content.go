package controllers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func DynamicContent(router fiber.Router) {

	content := router.Group("content")
	content.Get("/", handlers.GetContent)
	content.Get("/dictionary", handlers.GetDictionary)
	content.Get("/dynamic", handlers.GetDynamicContent)
}
