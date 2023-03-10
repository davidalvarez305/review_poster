package types

type DynamicContent struct {
	Sentences string `json:"sentences"`
	Template  string `json:"template"`
	Paragraph string `json:"paragraph"`
}

type Dictionary struct {
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Synonyms string `json:"synonyms"`
}

type ProcessedDictionary struct {
	Word     string   `json:"word"`
	Tag      string   `json:"tag"`
	Synonyms []string `json:"synonyms"`
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
