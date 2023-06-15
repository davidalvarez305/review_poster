package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func User(router fiber.Router) {

	user := router.Group("user", middleware.PostMiddleware)

	user.Get("/", middleware.ResourceAccessRestriction(handlers.GetUser))
	user.Put("/", middleware.ResourceAccessRestriction(handlers.UpdateUser))
	user.Delete("/", middleware.ResourceAccessRestriction(handlers.DeleteUser))
	user.Post("/register", middleware.ResourceAccessRestriction(handlers.CreateUser))
	user.Post("/login", middleware.ResourceAccessRestriction(handlers.Login))
	user.Post("/logout", middleware.ResourceAccessRestriction(handlers.Logout))
	user.Get("/change-password", middleware.ResourceAccessRestriction(handlers.RequestChangePasswordCode))
	user.Put("/change-password/:code", middleware.ResourceAccessRestriction(handlers.ChangePassword))

	// User Words Resources
	user.Get("/:userId/word", middleware.ResourceAccessRestriction(handlers.GetUserWords))
	user.Post("/:userId/word", middleware.ResourceAccessRestriction(handlers.CreateUserWord))
	user.Put("/:userId/word/:wordId", middleware.ResourceAccessRestriction(handlers.UpdateUserWord))

	// User Synonyms By Word
	user.Get("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.GetUserSynonymsByWord))
	user.Patch("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.UpdateUserSynonymsByWord))
	user.Post("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.CreateUserSynonymsByWord))
	user.Delete("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.DeleteUserSynonymsByWord))
	user.Put("/:userId/word/:word/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.UpdateUserSynonymByWord))
	user.Delete("/:userId/word/:word/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.DeleteUserSynonymByWord))

	// User Templates Resource
	user.Get("/:userId/template", middleware.ResourceAccessRestriction(handlers.GetUserTemplates))
	user.Post("/:userId/template", middleware.ResourceAccessRestriction(handlers.CreateUserTemplates))
	user.Put("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.UpdateUserTemplate))
	user.Delete("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.DeleteUserTemplate))
	user.Get("/:userId/template/:templateName/sentence", middleware.ResourceAccessRestriction(handlers.GetUserSentencesByTemplate))
	user.Get("/:userId/template/:templateName/paragraph/sentence", middleware.ResourceAccessRestriction(handlers.GetUserJoinedSentencesByParagraph))

	// User Paragraphs Resource
	user.Get("/:userId/template/:templateName/paragraph", middleware.ResourceAccessRestriction(handlers.GetUserParagraphsByTemplate))
	user.Post("/:userId/template/:templateName/paragraph", middleware.ResourceAccessRestriction(handlers.CreateUserParagraphsByTemplate))
	user.Patch("/:userId/template/:templateName/paragraph", middleware.ResourceAccessRestriction(handlers.UpdateUserParagraphsByTemplate))
	user.Delete("/:userId/template/:templateName/paragraph", middleware.ResourceAccessRestriction(handlers.DeleteUserParagraphsByTemplate))
	user.Put("/:userId/template/:templateName/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.UpdateUserParagraphByTemplate))
	user.Delete("/:userId/template/:templateName/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.DeleteUserParagraphByTemplate))

	// User Sentence By Paragraph Resource
	user.Get("/:userId/template/:templateName/paragraph/:paragraphName/sentence", middleware.ResourceAccessRestriction(handlers.GetUserParagraphSentencesByTemplate))
	user.Post("/:userId/template/:templateName/paragraph/:paragraphName/sentence", middleware.ResourceAccessRestriction(handlers.CreateUserParagraphSentencesByTemplate))
	user.Patch("/:userId/template/:templateName/paragraph/:paragraphName/sentence", middleware.ResourceAccessRestriction(handlers.UpdateUserParagraphSentencesByTemplate))
	user.Delete("/:userId/template/:templateName/paragraph/:paragraphName/sentence", middleware.ResourceAccessRestriction(handlers.DeleteUserParagraphSentencesByTemplate))
	user.Put("/:userId/template/:templateName/paragraph/:paragraphName/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.UpdateParagraphSentenceByTemplate))
	user.Delete("/:userId/template/:templateName/paragraph/:paragraphName/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.DeleteParagraphSentenceByTemplate))
}
