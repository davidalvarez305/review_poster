package models

type SubCategory struct {
	ID          int           `json:"id" form:"id"`
	Name        string        `gorm:"unique;column:name" json:"name" form:"name"`
	Slug        string        `gorm:"unique;column:slug" json:"slug" form:"slug"`
	ReviewPosts []*ReviewPost `json:"review_posts" form:"review_posts"`
	CategoryID  int           `json:"category_id" form:"category_id"`
	Category    *Category     `gorm:"not null;column:category_id;foreignKey:CategoryID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"category" form:"category"`
}
