package main

import (
	"encoding/gob"
	"fmt"
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

	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
		log.Fatalf("Error loading env file.: %+v\n", err)
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
