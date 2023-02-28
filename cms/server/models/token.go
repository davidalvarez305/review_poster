package models

type Token struct {
	ID        int    `json:"id" form:"id"`
	UUID      string `gorm:"unique;not null" json:"uuid" form:"uuid"`
	CreatedAt int64  `gorm:"not null;column:created_at" json:"created_at" form:"created_at"`
}
