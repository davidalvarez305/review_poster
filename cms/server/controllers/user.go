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
	user.Get("/word/:userId", handlers.GetWords)
	user.Get("/word/:userId/:word", handlers.GetWord)
	user.Post("/word/:userId", handlers.CreateWord)
	user.Put("/word/:userId/:id", handlers.UpdateWord)
	user.Delete("/word/:userId/:id", handlers.DeleteWord)
}
