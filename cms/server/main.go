package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
	"github.com/davidalvarez305/content_go/server/routes"
	"github.com/davidalvarez305/content_go/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	gob.Register(models.User{})

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading env file.")
	}

	CLIENT_URL := os.Getenv("CLIENT_URL")
	PORT := os.Getenv("PORT")

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
