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

	// Words related endpoints
	user.Get("/:userId/word", middleware.ResourceAccessRestriction(handlers.GetWords))
	user.Post("/:userId/word", middleware.ResourceAccessRestriction(handlers.CreateWord))
	user.Put("/:userId/word/:wordId", middleware.ResourceAccessRestriction(handlers.UpdateWord))
	user.Delete("/:userId/word/:wordId", middleware.ResourceAccessRestriction(handlers.DeleteWord))

	// User synonyms by word
	user.Get("/:userId/word/:word/synonym", middleware.ResourceAccessRestriction(handlers.GetUserSynonymsByWord))

	// Template related endpoints
	user.Get("/:userId/template", middleware.ResourceAccessRestriction(handlers.GetTemplates))
	user.Post("/:userId/template", middleware.ResourceAccessRestriction(handlers.CreateTemplate))
	user.Put("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.UpdateTemplate))
	user.Delete("/:userId/template/:templateId", middleware.ResourceAccessRestriction(handlers.DeleteTemplate))

	// Content related endpoints
	user.Get("/:userId/content", middleware.ResourceAccessRestriction(handlers.GetContent))

	// Paragraph related endpoints
	user.Get("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.GetParagraphs))
	user.Post("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.CreateParagraphs))
	user.Put("/:userId/paragraph", middleware.ResourceAccessRestriction(handlers.UpdateParagraphs))
	user.Put("/:userId/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.UpdateParagraph))
	user.Delete("/:userId/paragraph/:paragraphId", middleware.ResourceAccessRestriction(handlers.DeleteParagraph))
	user.Post("/:userId/paragraph/bulk", middleware.ResourceAccessRestriction(handlers.BulkParagraphsUpdate))

	// Sentence related endpoints
	user.Get("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.GetSentences))
	user.Post("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.CreateSentences))
	user.Put("/:userId/sentence", middleware.ResourceAccessRestriction(handlers.UpdateSentences))
	user.Put("/:userId/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.UpdateSentence))
	user.Delete("/:userId/sentence/:sentenceId", middleware.ResourceAccessRestriction(handlers.DeleteSentence))
	user.Post("/:userId/sentence/bulk", middleware.ResourceAccessRestriction(handlers.BulkSentencesUpdate))

	// Synonym related endpoints
	user.Get("/:userId/synonym", middleware.ResourceAccessRestriction(handlers.GetSelectedSynonyms))
	user.Post("/:userId/synonym", middleware.ResourceAccessRestriction(handlers.CreateSynonym))
	user.Put("/:userId/synonym", middleware.ResourceAccessRestriction(handlers.UpdateSynonyms))
	user.Put("/:userId/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.UpdateSynonym))
	user.Delete("/:userId/synonym/:synonymId", middleware.ResourceAccessRestriction(handlers.DeleteSynonym))
	user.Post("/:userId/synonym/bulk", middleware.ResourceAccessRestriction(handlers.BulkSynonymsPost))
}
