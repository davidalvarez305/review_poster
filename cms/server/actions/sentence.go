package actions

import (
	"fmt"

	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
)

type Sentence struct {
	*models.Sentence
}

type Sentences []*models.Sentence

type JoinedParagraph struct {
	ParagraphID    int    `json:"paragraph_id"`
	ParagraphName  string `json:"paragraph_name"`
	ParagraphOrder int    `json:"paragraph_order"`
	TemplateID     uint   `json:"template_id"`
	TemplateName   string `json:"template_name"`
}

type JoinedParagraphs []*JoinedParagraph

func (sentences *Sentences) GetSentencesByParagraph(paragraph, userId string) error {
	sql := fmt.Sprintf(`SELECT * FROM sentence WHERE paragraph_id = (
		SELECT id FROM paragraph WHERE name = '%s' AND user_id = '%s'
	)`, paragraph, userId)
	result := database.DB.Raw(sql).Scan(&sentences)
	return result.Error
}

func (sentences *Sentences) GetAllSentences(userId string) error {
	result := database.DB.Where("user_id = ?", userId).Find(&sentences)
	return result.Error
}

// Create a single sentence. This assumes that input will have been validated elsewhere.
func (sentence *Sentence) CreateSentence() error {
	result := database.DB.Save(&sentence).First(&sentence)
	return result.Error
}

// Create many sentences. This assumes that input will have been validated elsewhere.
func (sentences *Sentences) CreateSentences() error {
	result := database.DB.Save(&sentences).Find(&sentences)
	return result.Error
}

func (sentences *Sentences) UpdateSentences(paragraph, userId string) error {
	result := database.DB.Where("user_id = ?", userId).Save(&sentences)

	if result.Error != nil {
		return result.Error
	}

	err := sentences.GetSentencesByParagraph(paragraph, userId)

	return err
}

func (sentences *Sentences) DeleteSentences(ids []int, paragraph, userId string) error {
	res := database.DB.Delete(&models.Sentence{}, ids)

	if res.Error != nil {
		return res.Error
	}

	err := sentences.GetSentencesByParagraph(paragraph, userId)

	return err
}

// Delete records without checking user or paragraph id. Assumes that this is being checked elsewhere.
func (sentences *Sentences) SimpleDelete() error {
	res := database.DB.Delete(&sentences)
	return res.Error
}

func (jp *JoinedParagraphs) GetTemplatesAndParagraphs(userId string) error {
	sql := fmt.Sprintf(`SELECT
	p.id AS paragraph_id, p.name AS paragraph_name, p.order AS paragraph_order, p.template_id,
	t.name AS template_name
	FROM paragraph AS p
	JOIN template AS t
	ON t.id = p.template_id
	WHERE t.user_id = %s;
	`, userId)
	res := database.DB.Raw(sql).Scan(&jp)
	return res.Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func (sentences Sentences) DeleteBulkSentences(clientSentences, existingSentences *Sentences) error {
	for _, existingSentence := range *existingSentences {
		found := false
		for _, clientSentence := range *clientSentences {
			if existingSentence.Sentence == clientSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, existingSentence)
		}
	}

	if len(sentences) > 0 {
		err := sentences.SimpleDelete()

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func (sentences Sentences) AddBulkSentences(clientSentences, existingSentences *Sentences) error {
	for _, clientSentence := range *clientSentences {
		found := false
		for _, existingSentence := range *existingSentences {
			if clientSentence.Sentence == existingSentence.Sentence {
				found = true
			}
		}
		if !found {
			sentences = append(sentences, clientSentence)
		}
	}

	if len(sentences) > 0 {
		err := sentences.CreateSentences()

		if err != nil {
			return err
		}
	}

	return nil
}
