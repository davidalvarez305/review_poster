package types

import "github.com/davidalvarez305/review_poster/server/models"

type CreatePostInput struct {
	ProductName    string `json:"productName"`
	Category       string `json:"category"`
	ProductURL     string `json:"productUrl"`
	ImageURL       string `json:"imageUrl"`
	ProductReviews string `json:"productReviews"`
	ProductPrice   string `json:"productPrice"`
	ProductRating  string `json:"productRating"`
}

type Keyword struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type CreateReviewPostsInput struct {
	Keyword   string `json:"keyword" form:"keyword"`
	GroupName string `json:"group_name" form:"group_name"`
	Template  string `json:"template" form:"template"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Prompt struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

type CreateWordInput struct {
	ID       int      `json:"id"`
	Word     string   `json:"word"`
	Tag      string   `json:"tag"`
	Synonyms []string `json:"synonyms"`
}

type UpdateUserSynonymsByWordInput struct {
	Synonyms       []models.Synonym `json:"synonyms" form:"synonyms"`
	DeleteSynonyms []int            `json:"delete_synonyms" form:"delete_synonyms"`
}
