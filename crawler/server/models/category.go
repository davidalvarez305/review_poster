package models

type Category struct {
	ID            int            `json:"id" form:"id"`
	Name          string         `gorm:"unique;column:name" json:"name" form:"name"`
	Slug          string         `gorm:"unique;column:slug" json:"slug" form:"slug"`
	SubCategories []*SubCategory `json:"sub_categories" form:"sub_categories"`
	GroupID       int            `json:"group_id" form:"group_id"`
	Group         *Group         `gorm:"not null;column:group_id;foreignKey:GroupID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"group" form:"group"`
}
