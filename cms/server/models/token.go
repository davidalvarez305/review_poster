package models

type Token struct {
	ID        uint   `json:"-"`
	UUID      string `gorm:"unique;not null" json:"uuid"`
	UserID    uint   `json:"user_id"`
	User      *Users `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	CreatedAt int64  `gorm:"not null;column:created_at" json:"createdAt"`
}
