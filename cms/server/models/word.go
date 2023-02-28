package models

type Word struct {
	ID       int        `json:"id" form:"id"`
	Name     string     `gorm:"unique;not null" json:"name" form:"name"`
	Tag      string     `gorm:"unique;not null" json:"tag" form:"tag"`
	UserID   int        `json:"user_id" form:"user_id"`
	User     *User      `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user" form:"user"`
	Synonyms []*Synonym `json:"synonyms" form:"synonyms"`
}
