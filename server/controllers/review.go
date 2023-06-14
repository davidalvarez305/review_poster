package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func ReviewPost(router fiber.Router) {

	reviewPost := router.Group("review-post", middleware.AuthMiddleware)

	reviewPost.Post("/", handlers.CreatePosts)
}
