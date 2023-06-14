package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func Amazon(router fiber.Router) {

	amazon := router.Group("amazon", middleware.AuthMiddleware)

	amazon.Get("/crawl", handlers.CrawlAmazon)
	amazon.Get("/paapi5", handlers.SearchPAAPI5)
}
