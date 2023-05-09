package models

type Template struct {
	ID     int    `json:"id" form:"id"`
	Name   string `gorm:"unique;not null" json:"name" form:"name"`
	UserID int    `json:"user_id" form:"user_id"`
	User   *User  `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user" form:"user"`
}
