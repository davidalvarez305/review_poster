package models

type Category struct {
	*Base
	SubCategories []*SubCategory `json:"sub_categories" form:"sub_categories"`
	GroupID       int            `json:"group_id" form:"group_id"`
	Group         *Group         `gorm:"not null;column:group_id;foreignKey:GroupID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"group" form:"group"`
}
