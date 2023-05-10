package server

import (
	"os"

	"github.com/davidalvarez305/review_poster/server/controllers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type Server struct {
	App     *fiber.App
	DB      *gorm.DB
	Session *session.Store
	Port    string
}

func NewServer(opts *Server) *Server {
	return &Server{
		App:     opts.App,
		DB:      opts.DB,
		Session: opts.Session,
		Port:    opts.Port,
	}
}

func (server *Server) Start() {

	CLIENT_URL := os.Getenv("CONTENT_CLIENT_URL")

	server.App.Use(cors.New(cors.Config{
		AllowOrigins:     CLIENT_URL,
		AllowCredentials: true,
	}))

	api := server.App.Group("api", middleware.AuthMiddleware, middleware.ResourceAccessRestriction)

	controllers.Google(api)
	controllers.Amazon(api)
	controllers.ReviewPost(api)
	controllers.DynamicContent(api)

	server.App.Listen(":" + server.Port)
}
