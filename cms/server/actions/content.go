package actions

import (
	"fmt"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Sentences []*models.Sentence

type Words []*models.Word

func (s *Sentences) GetSentences(template string) error {
	return database.DB.Where("template_id = ?", template).Preload("Sentence.Template").Preload("Sentence.Paragraph").Find(&s).Error
}

func (w *Words) GetWords() error {
	return database.DB.Preload("Word.Synonym").Find(&w).Error
}
