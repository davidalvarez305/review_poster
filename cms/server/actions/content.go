package actions

import (
	"fmt"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Sentences []*models.Sentence

type Words []*models.Word

func (s *Sentences) GetSentences(template, userId string) error {
	return database.DB.Where("template_id = ? AND user_id = ?", template, userId).Preload("Sentence.Template").Preload("Sentence.Paragraph").Find(&s).Error
}

func (w *Words) GetWords(userId string) error {
	return database.DB.Where("user_id = ?", userId).Preload("Word.Synonym").Find(&w).Error
}
