package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/routes"
	"github.com/davidalvarez305/review_poster/cms/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	gob.Register(models.User{})

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading env file.: %+v\n", err)
	}

	CLIENT_URL := os.Getenv("CONTENT_CLIENT_URL")
	PORT := os.Getenv("CMS_PORT")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     CLIENT_URL,
		AllowCredentials: true,
	}))

	database.Connect()
	sessions.Init()

	routes.Router(app)
	app.Listen(":" + PORT)
}
