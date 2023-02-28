package models

type Synonym struct {
	ID      int    `json:"id" form:"id"`
	Synonym string `gorm:"unique;not null" json:"synonym" form:"synonym"`
	WordID  int    `json:"word_id" form:"word_id"`
	Word    *Word  `gorm:"not null;column:word_id;foreignKey:WordID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"word" form:"word"`
}
