package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Word(router fiber.Router) {
	word := router.Group("word")

	word.Get("/:userId", handlers.GetWords)
	word.Get("/:userId/:word", handlers.GetWord)
	word.Post("/:userId", handlers.CreateWord)
	word.Put("/:userId/:id", handlers.UpdateWord)
	word.Delete("/:userId/:id", handlers.DeleteWord)

}
