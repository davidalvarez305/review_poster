package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Template(router fiber.Router) {
	template := router.Group("template")

	template.Get("/", handlers.GetTemplates)
	template.Post("/", handlers.CreateTemplate)
	template.Put("/", handlers.UpdateTemplate)
	template.Delete("/", handlers.DeleteTemplate)

}
