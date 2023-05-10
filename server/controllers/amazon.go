package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Amazon(router fiber.Router) {

	amazon := router.Group("amazon")

	amazon.Get("/crawl", handlers.CrawlAmazon)
	amazon.Get("/paapi5", handlers.SearchPAAPI5)
}
