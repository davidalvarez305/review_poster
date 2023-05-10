package types

import (
	"github.com/davidalvarez305/review_poster/server/models"
)

type DictionaryAPIResponse struct {
	Data []models.Word `json:"data"`
}

type ContentAPIResponse struct {
	Data []models.Sentence `json:"data"`
}

type Synonym struct {
	ID      int          `json:"id" form:"id"`
	Synonym string       `json:"synonym" form:"synonym"`
	WordID  int          `json:"word_id" form:"word_id"`
	Word    *models.Word `json:"word" form:"word"`
}

type ProcessedContent struct {
	ReviewPostTitle              []string `json:"ReviewPostTitle"`
	ReviewPostContent            []string `json:"ReviewPostContent"`
	ReviewPostHeadline           []string `json:"ReviewPostHeadline"`
	ReviewPostIntro              []string `json:"ReviewPostIntro"`
	ReviewPostDescription        []string `json:"ReviewPostDescription"`
	ReviewPostProductLabel       []string `json:"ReviewPostProductLabel"`
	ReviewPostProductDescription []string `json:"ReviewPostProductDescription"`
	ReviewPostFaq_Answer_1       []string `json:"ReviewPostFaq_Answer_1"`
	ReviewPostFaq_Answer_2       []string `json:"ReviewPostFaq_Answer_2"`
	ReviewPostFaq_Answer_3       []string `json:"ReviewPostFaq_Answer_3"`
	ReviewPostFaq_Question_1     []string `json:"ReviewPostFaq_Question_1"`
	ReviewPostFaq_Question_2     []string `json:"ReviewPostFaq_Question_2"`
	ReviewPostFaq_Question_3     []string `json:"ReviewPostFaq_Question_3"`
}

type FinalizedContent struct {
	ReviewPostTitle              string `json:"ReviewPostTitle"`
	ReviewPostContent            string `json:"ReviewPostContent"`
	ReviewPostHeadline           string `json:"ReviewPostHeadline"`
	ReviewPostIntro              string `json:"ReviewPostIntro"`
	ReviewPostDescription        string `json:"ReviewPostDescription"`
	ReviewPostProductLabel       string `json:"ReviewPostProductLabel"`
	ReviewPostProductDescription string `json:"ReviewPostProductDescription"`
	ReviewPostFaq_Answer_1       string `json:"ReviewPostFaq_Answer_1"`
	ReviewPostFaq_Answer_2       string `json:"ReviewPostFaq_Answer_2"`
	ReviewPostFaq_Answer_3       string `json:"ReviewPostFaq_Answer_3"`
	ReviewPostFaq_Question_1     string `json:"ReviewPostFaq_Question_1"`
	ReviewPostFaq_Question_2     string `json:"ReviewPostFaq_Question_2"`
	ReviewPostFaq_Question_3     string `json:"ReviewPostFaq_Question_3"`
}
