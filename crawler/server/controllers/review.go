package controllers

import (
	"github.com/davidalvarez305/review_poster/crawler/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func ReviewPost(router fiber.Router) {

	reviewPost := router.Group("review-post")

	reviewPost.Post("/", handlers.CreatePosts)
}
