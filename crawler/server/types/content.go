package types

type DictionaryAPIResponse struct {
	Data []Word `json:"data"`
}

type ContentAPIResponse struct {
	Data []Sentence `json:"data"`
}

type Synonym struct {
	ID      int    `json:"id" form:"id"`
	Synonym string `json:"synonym" form:"synonym"`
	WordID  int    `json:"word_id" form:"word_id"`
	Word    *Word  `json:"word" form:"word"`
}

type Word struct {
	ID       int        `json:"id" form:"id"`
	Name     string     `json:"name" form:"name"`
	Tag      string     `json:"tag" form:"tag"`
	UserID   int        `json:"user_id" form:"user_id"`
	Synonyms []*Synonym `json:"synonyms" form:"synonyms"`
}

type Sentence struct {
	ID          int        `json:"id" form:"id"`
	Sentence    string     `gorm:"unique;not null" json:"sentence" form:"sentence"`
	ParagraphID int        `json:"paragraph_id" form:"paragraph_id"`
	Paragraph   *Paragraph `json:"paragraph" form:"paragraph"`
	TemplateID  int        `json:"template_id" form:"template_id"`
	Template    *Template  `json:"template" form:"template"`
}

type Paragraph struct {
	ID         int       `json:"id" form:"id"`
	Name       string    `gorm:"unique;not null" json:"name" form:"name"`
	Order      int       `gorm:"null:true;default:null" json:"order" form:"order"`
	TemplateID int       `json:"template_id" form:"template_id"`
	Template   *Template `json:"template" form:"template"`
	UserID     int       `json:"user_id" form:"user_id"`
}

type Template struct {
	ID     int    `json:"id" form:"id"`
	Name   string `gorm:"unique;not null" json:"name" form:"name"`
	UserID int    `json:"user_id" form:"user_id"`
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
