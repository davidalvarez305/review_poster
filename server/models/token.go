package models

type Token struct {
	ID        int    `json:"id" form:"id"`
	UUID      string `gorm:"unique;not null" json:"uuid" form:"uuid"`
	CreatedAt int64  `gorm:"not null;column:created_at" json:"created_at" form:"created_at"`
	UserID    int    `json:"user_id" form:"user_id"`
	User      *User  `gorm:"unique;not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user" form:"user"`
}
