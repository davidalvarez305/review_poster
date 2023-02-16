package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Word(router fiber.Router) {
	word := router.Group("word")

	word.Get("/", handlers.GetWords)
	word.Get("/:word", handlers.GetWord)
	word.Post("/", handlers.CreateWord)
	word.Put("/", handlers.UpdateWord)
	word.Delete("/", handlers.DeleteWord)

}
