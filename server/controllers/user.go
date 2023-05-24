package controllers

import (
	"github.com/davidalvarez305/review_poster/server/handlers"
	"github.com/davidalvarez305/review_poster/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func User(router fiber.Router) {

	user := router.Group("user")

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

	// User Word Resource
	user.Get("/:userId/word/:wordName", middleware.ResourceAccessRestriction(handlers.GetUserWord))
	user.Put("/:userId/word/:wordId", middleware.ResourceAccessRestriction(handlers.UpdateUserWord))
	user.Delete("/:userId/word/:wordId", middleware.ResourceAccessRestriction(handlers.DeleteUserWord))

	// User Synonyms By Word Resources
	user.Get("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.GetUserSynonymsByWord))
	user.Put("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.UpdateUserSynonymsByWord))
	user.Post("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.CreateUserSynonymsByWord))

	// User Templates Resource
	user.Get("/:userId/template", middleware.ResourceAccessRestriction(handlers.GetTemplates))
	user.Post("/:userId/template", middleware.ResourceAccessRestriction(handlers.CreateTemplate))

	// User Template Resource
	user.Put("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.UpdateTemplate))
	user.Delete("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.DeleteTemplate))

	// Content related endpoints
	user.Get("/:userId/content", middleware.ResourceAccessRestriction(handlers.GetContent))

	// User Paragraphs Resource
	user.Get("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.GetParagraphs))
	user.Post("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.CreateParagraphs))
	user.Put("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.UpdateParagraphs))

	// User Paragraph Resource
	user.Put("/:userId/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.UpdateParagraph))
	user.Delete("/:userId/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.DeleteParagraph))
	user.Post("/:userId/paragraph/bulk", middleware.ResourceAccessRestriction(handlers.BulkParagraphsUpdate))

	// User Sentence By Paragraph Resource
	user.Get("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.GetSentences))
	user.Post("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.CreateSentences))
	user.Put("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.UpdateSentences))

	// User Sentences By Paragraph Resource
	user.Put("/:userId/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.UpdateSentence))
	user.Delete("/:userId/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.DeleteSentence))
	user.Post("/:userId/sentence/bulk", middleware.ResourceAccessRestriction(handlers.BulkSentencesUpdate))

	// User Synonyms By Word Resource
	user.Post("/:userId/synonym", middleware.ResourceAccessRestriction(handlers.CreateSynonym))
	user.Put("/:userId/synonym", middleware.ResourceAccessRestriction(handlers.UpdateSynonyms))

	// User Synonym By Word Resource
	user.Put("/:userId/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.UpdateSynonym))
	user.Delete("/:userId/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.DeleteSynonym))
}
