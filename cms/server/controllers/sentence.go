package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Sentence(router fiber.Router) {
	sentence := router.Group("sentence")

	sentence.Get("/", handlers.GetSentences)
	sentence.Get("/create-sentence", handlers.GetTemplatesAndParagraphs)
	sentence.Get("/:paragraph", handlers.GetSentencesByParagraph)
	sentence.Post("/", handlers.CreateSentences)
	sentence.Put("/", handlers.UpdateSentences)
	sentence.Delete("/", handlers.DeleteSentence)
	sentence.Post("/bulk", handlers.BulkSentencesUpdate)
}
