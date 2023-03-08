package models

type User struct {
	ID       int    `json:"id" form:"id"`
	Username string `gorm:"unique;not null" json:"username" form:"username"`
	Password string `gorm:"not null" json:"password" form:"password"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	TokenID  int    `json:"token_id" form:"token_id"`
	Token    *Token `gorm:"unique;not null;column:token;foreignKey:TokenID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"token" form:"token"`
}
