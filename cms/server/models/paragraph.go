package models

type Paragraph struct {
	ID         int       `json:"id" form:"id"`
	Name       string    `gorm:"unique;not null" json:"name" form:"name"`
	Order      int       `gorm:"null:true;default:null" json:"order" form:"order"`
	TemplateID int       `json:"template_id" form:"template_id"`
	Template   *Template `gorm:"not null;column:template_id;foreignKey:TemplateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"template" form:"template"`
	UserID     int       `json:"user_id" form:"user_id"`
	User       *User     `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"user" form:"user"`
}
