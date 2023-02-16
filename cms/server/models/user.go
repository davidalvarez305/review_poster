package models

type Users struct {
	ID       uint   `json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"unique;not null" json:"email"`
	Token    string `gorm:"unique;not null" json:"token"`
}
