package models

type Word struct {
	ID     int    `json:"id"`
	Name   string `gorm:"unique;not null" json:"name"`
	Tag    string `gorm:"unique;not null" json:"tag"`
	UserID uint   `json:"user_id"`
	User   *Users `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}
