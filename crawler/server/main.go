package main

import (
	"log"
	"os"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/routes"
	"github.com/davidalvarez305/review_poster/crawler/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("ERROR LOADING ENV FILE: %+v\n", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	err = database.Connect()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}

	// Can't call this and return error, it will panic on its own...
	sessions.Init()

	routes.Router(app)

	app.Listen(":" + os.Getenv("CRAWLER_PORT"))
}
