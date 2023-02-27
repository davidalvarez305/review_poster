package models

type SubCategory struct {
	*Base
	ReviewPosts []*ReviewPost `json:"review_posts" form:"review_posts"`
	CategoryID  int           `json:"category_id" form:"category_id"`
	Category    *Category     `gorm:"not null;column:category_id;foreignKey:CategoryID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"category" form:"category"`
}
