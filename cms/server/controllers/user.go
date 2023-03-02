package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func User(router fiber.Router) {

	user := router.Group("user")

	user.Get("/", handlers.GetUser)
	user.Put("/", handlers.UpdateUser)
	user.Delete("/", handlers.DeleteUser)
	user.Post("/register", handlers.CreateUser)
	user.Post("/login", handlers.Login)
	user.Post("/logout", handlers.Logout)
	user.Get("/change-password", handlers.RequestChangePasswordCode)
	user.Put("/change-password/:code", handlers.ChangePassword)

	// Words related endpoints
	user.Get("/:userId/word", handlers.GetWords)
	user.Get("/:userId/word/:word", handlers.GetWord)
	user.Post("/:userId/word", handlers.CreateWord)
	user.Put("/:userId/word/:id", handlers.UpdateWord)
	user.Delete("/:userId/word/:id", handlers.DeleteWord)

	// Template related endpoints
	user.Get("/:userId/template", handlers.GetTemplates)
	user.Post("/:userId/template", handlers.CreateTemplate)
	user.Put("/:userId/template/:templateId", handlers.UpdateTemplate)
	user.Delete("/:userId/template/:templateId", handlers.DeleteTemplate)
}
