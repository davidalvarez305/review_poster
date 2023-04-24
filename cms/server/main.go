package main

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/server"
	"github.com/davidalvarez305/review_poster/cms/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	gob.Register(models.User{})

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading env file.: %+v\n", err)
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}

	sessionStore := sessions.Init()

	server := server.NewServer(&server.Server{
		App:   fiber.New(),
		DB:    db,
		Store: sessionStore,
		Port:  os.Getenv("CMS_PORT"),
	})

	server.Start()
}
