package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Paragraph(router fiber.Router) {
	paragraph := router.Group("paragraph")

	paragraph.Get("/", handlers.GetParagraphs)
	paragraph.Post("/", handlers.CreateParagraphs)
	paragraph.Put("/", handlers.UpdateParagraphs)
	paragraph.Delete("/", handlers.DeleteParagraph)
	paragraph.Get("/selected/", handlers.GetSelectedParagraphs)
	paragraph.Post("/bulk", handlers.BulkParagraphsUpdate)
}
