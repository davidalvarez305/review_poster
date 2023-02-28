package models

type Sentence struct {
	ID          int        `json:"id" form:"id"`
	Sentence    string     `gorm:"unique;not null" json:"sentence" form:"sentence"`
	ParagraphID int        `json:"paragraph_id" form:"paragraph_id"`
	Paragraph   *Paragraph `gorm:"not null;column:paragraph_id;foreignKey:ParagraphID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"paragraph" form:"paragraph"`
	TemplateID  int        `json:"template_id" form:"template_id"`
	Template    *Template  `gorm:"not null;column:template_id;foreignKey:TemplateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"template" form:"template"`
}
