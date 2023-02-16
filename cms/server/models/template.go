package models

type Template struct {
	ID     int    `json:"id"`
	Name   string `gorm:"unique;not null" json:"name"`
	UserID uint   `json:"user_id"`
	User   *Users `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}
