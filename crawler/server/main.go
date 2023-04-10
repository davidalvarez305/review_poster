package main

import (
	"fmt"
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
	err := godotenv.Load("../../.env")

	if err != nil {
		fmt.Printf("ERROR: %+v\n", err)
		log.Fatalf("Error loading env file.: %+v\n", err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	database.Connect()
	sessions.Init()
	routes.Router(app)

	app.Listen(":" + os.Getenv("PORT"))
}
