package models

type Paragraph struct {
	ID         int       `json:"id"`
	Name       string    `gorm:"unique;not null" json:"name"`
	Order      int       `gorm:"null:true;default:null" json:"order"`
	TemplateID uint      `json:"template_id"`
	Template   *Template `gorm:"not null;column:template_id;foreignKey:TemplateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	UserID     uint      `json:"user_id"`
	User       *Users    `gorm:"not null;column:user_id;foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}
