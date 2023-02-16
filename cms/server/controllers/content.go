package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Content(router fiber.Router) {
	content := router.Group("content")

	content.Get("/content", handlers.GetContent)
	content.Get("/dictionary", handlers.GetDictionary)

}
