package models

type Category struct {
	ID             int
	Name           string
	Slug           string `gorm:"unique"`
	ReviewProducts []ReviewPost
	CategoryGroup  CategoryGroup `gorm:"column:category_group_id"`
}
