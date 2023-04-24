package server

import (
	"os"

	"github.com/davidalvarez305/review_poster/cms/server/controllers"
	"github.com/davidalvarez305/review_poster/cms/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type Server struct {
	App   *fiber.App
	DB    *gorm.DB
	Store *session.Store
	Port  string
}

var DB gorm.DB

var Sessions session.Store

func NewServer(opts *Server) *Server {
	return &Server{
		App:   opts.App,
		DB:    opts.DB,
		Store: opts.Store,
		Port:  opts.Port,
	}
}

func (server *Server) Start() {

	CLIENT_URL := os.Getenv("CONTENT_CLIENT_URL")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     CLIENT_URL,
		AllowCredentials: true,
	}))

	api := server.App.Group("api", middleware.AuthMiddleware, middleware.ResourceAccessRestriction)

	controllers.User(api)

	DB = *server.DB

	server.App.Listen(":" + server.Port)
}
