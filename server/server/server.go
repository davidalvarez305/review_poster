package server

import (
	"github.com/davidalvarez305/review_poster/crawler/server/controllers"
	"github.com/davidalvarez305/review_poster/crawler/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type Server struct {
	App  *fiber.App
	DB   *gorm.DB
	Port string
}

func NewServer(opts *Server) *Server {
	return &Server{
		App:  opts.App,
		DB:   opts.DB,
		Port: opts.Port,
	}
}

func (server *Server) Start() {

	CLIENT_URL := os.Getenv("CONTENT_CLIENT_URL")

	server.App.Use(cors.New(cors.Config{
		AllowOrigins:     CLIENT_URL,
		AllowCredentials: true,
	}))

	api := app.Group("api", middleware.AuthMiddleware, middleware.ResourceAccessRestriction)

	controllers.Google(api)
	controllers.Amazon(api)
	controllers.ReviewPost(api)
	controllers.DynamicContent(api)

	server.App.Listen(":" + server.Port)
}