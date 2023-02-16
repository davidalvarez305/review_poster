package models

type Sentence struct {
	ID          int        `json:"id"`
	Sentence    string     `gorm:"unique;not null" json:"sentence"`
	ParagraphID int        `json:"paragraph_id"`
	Paragraph   *Paragraph `gorm:"not null;column:paragraph_id;foreignKey:ParagraphID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	TemplateID  int        `json:"template_id"`
	Template    *Template  `gorm:"not null;column:template_id;foreignKey:TemplateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}
