package models

type Base struct {
	ID   int    `json:"id" form:"id"`
	Name string `gorm:"unique;column:name" json:"name" form:"name"`
	Slug string `gorm:"unique;column:slug" json:"slug" form:"slug"`
}
