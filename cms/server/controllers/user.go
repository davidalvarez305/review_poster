package controllers

import (
	"github.com/davidalvarez305/content_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func User(router fiber.Router) {

	user := router.Group("user")

	user.Get("/", handlers.GetUser)
	user.Put("/", handlers.UpdateUser)
	user.Delete("/", handlers.DeleteUser)
	user.Post("/register", handlers.CreateUser)
	user.Post("/login", handlers.Login)
	user.Post("/logout", handlers.Logout)
	user.Get("/change-password", handlers.RequestChangePasswordCode)
	user.Put("/change-password/:code", handlers.ChangePassword)

	// Words related endpoints
	user.Get("/:userId/word", handlers.GetWords)
	user.Get("/:userId/word/:word", handlers.GetWord)
	user.Post("/:userId/word", handlers.CreateWord)
	user.Put("/:userId/word/:id", handlers.UpdateWord)
	user.Delete("/:userId/word/:id", handlers.DeleteWord)

	// Template related endpoints
	user.Get("/:userId/template", handlers.GetTemplates)
	user.Post("/:userId/template", handlers.CreateTemplate)
	user.Put("/:userId/template/:templateId", handlers.UpdateTemplate)
	user.Delete("/:userId/template/:templateId", handlers.DeleteTemplate)

	// Content related endpoints
	user.Get("/:userId/content", handlers.GetContent)
	user.Get("/:userId/dictionary", handlers.GetDictionary)

	// Paragraph related endpoints
	user.Get("/:userId/paragraph", handlers.GetParagraphs)
	user.Post("/:userId/paragraph", handlers.CreateParagraphs)
	user.Put("/:userId/paragraph/:paragraphId", handlers.UpdateParagraphs)
	user.Delete("/:userId/paragraph/:paragraphId", handlers.DeleteParagraph)
	user.Get("/:userId/paragraph/selected", handlers.GetSelectedParagraphs)
	user.Post("/:userId/paragraph/bulk", handlers.BulkParagraphsUpdate)

	// Sentence related endpoints
	user.Get("/:userId/sentence", handlers.GetSentences)
	user.Get("/:userId/sentence/:paragraph", handlers.GetSentencesByParagraph)
	user.Post("/:userId/sentence", handlers.CreateSentences)
	user.Put("/:userId/sentence", handlers.UpdateSentences)
	user.Delete("/:userId/sentence", handlers.DeleteSentence)
	user.Post("/:userId/sentence/bulk", handlers.BulkSentencesUpdate)

	// Synonym related endpoints
	user.Get("/:userId/synonym", handlers.GetSelectedSynonyms)
	user.Post("/:userId/synonym", handlers.CreateSynonym)
	user.Put("/:userId/synonym", handlers.UpdateSynonyms)
	user.Delete("/:userId/synonym", handlers.DeleteSynonym)
	user.Post("/:userId/synonym/bulk", handlers.BulkSynonymsPost)
}
