package models

type Group struct {
	ID         int         `json:"id" form:"id"`
	Name       string      `gorm:"unique;column:name" json:"name" form:"name"`
	Slug       string      `gorm:"unique;column:slug" json:"slug" form:"slug"`
	Categories []*Category `json:"categories" form:"categories"`
}
