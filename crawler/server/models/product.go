package models

type Product struct {
	ID             int    `json:"id" form:"id"`
	AffiliateUrl   string `gorm:"unique;column:affiliate_url" json:"affiliate_url" form:"affiliate_url"`
	ProductPrice   string `gorm:"column:product_price" json:"product_price" form:"product_price"`
	ProductReviews string `gorm:"column:product_reviews" json:"product_reviews" form:"product_reviews"`
	ProductRatings string `gorm:"column:product_ratings" json:"product_ratings" form:"product_ratings"`
	ProductImage   string `gorm:"column:product_image" json:"product_image" form:"product_image"`
}
