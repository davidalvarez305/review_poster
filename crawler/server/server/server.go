package server

import (
	"github.com/davidalvarez305/review_poster/crawler/server/controllers"
	"github.com/davidalvarez305/review_poster/crawler/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	DB := &gorm.DB{}

	server.App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	api := server.App.Group("api", middleware.AuthMiddleware)

	controllers.Google(api)
	controllers.Amazon(api)
	controllers.ReviewPost(api)
	controllers.DynamicContent(api)

	DB = server.DB

	server.App.Listen(":" + server.Port)
}
