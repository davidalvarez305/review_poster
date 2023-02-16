package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Synonym(router fiber.Router) {
	synonym := router.Group("synonym")

	synonym.Get("/", handlers.GetSelectedSynonyms)
	synonym.Post("/", handlers.CreateSynonym)
	synonym.Put("/", handlers.UpdateSynonyms)
	synonym.Delete("/", handlers.DeleteSynonym)
	synonym.Post("/bulk", handlers.BulkSynonymsPost)
}
