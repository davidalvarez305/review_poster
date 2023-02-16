package sessions

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/postgres"
)

var Sessions *session.Store

func Init() {
	storage := postgres.New(postgres.Config{
		Username:   os.Getenv("POSTGRES_USER"),
		Password:   os.Getenv("POSTGRES_PASSWORD"),
		Host:       os.Getenv("POSTGRES_HOST"),
		Port:       5432,
		Database:   "fiber",
		Table:      "fiber_storage",
		SslMode:    "disable",
		GCInterval: 60 * 60 * 24 * 365 * time.Second,
	})

	key := fmt.Sprintf("cookie:%s", os.Getenv("COOKIE_NAME"))

	store := session.New(session.Config{
		Expiration:     24 * 365 * time.Hour,
		Storage:        storage,
		KeyLookup:      key,
		CookieSameSite: "lax",
		KeyGenerator:   utils.UUID,
	})

	Sessions = store
	fmt.Println(store)
}
