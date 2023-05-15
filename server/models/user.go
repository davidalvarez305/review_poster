package models

type User struct {
	ID               int    `json:"id" form:"id"`
	Username         string `gorm:"unique;not null" json:"username" form:"username"`
	Password         string `gorm:"not null" json:"password" form:"password"`
	Email            string `gorm:"unique;not null" json:"email" form:"email"`
	AuthHeaderString string `json:"auth_header_string" form:"auth_header_string"`
}
