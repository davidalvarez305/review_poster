package main

import (
	"log"
	"os"

	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/server"
	"github.com/davidalvarez305/review_poster/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("ERROR LOADING ENV FILE: %+v\n", err)
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}

	sessionStore := sessions.Init()

	if err != nil {
		log.Fatalf("ERROR CONNECTING TO DB: %+v\n", err)
	}

	server := server.NewServer(&server.Server{
		App:     fiber.New(),
		DB:      db,
		Session: sessionStore,
		Port:    os.Getenv("SERVER_PORT"),
	})

	server.Start()
}
