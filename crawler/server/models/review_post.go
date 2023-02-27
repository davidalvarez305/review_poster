package models

type ReviewPost struct {
	ID                            uint
	Title                         string `json:"title" form:"title"`
	Slug                          string `gorm:"unique" json:"slug" form:"slug"`
	Content                       string `json:"content" form:"content"`
	CategoryID                    int    `gorm:"column:category_id" json:"category_id" form:"category_id"`
	Category                      *Category
	Headline                      string `json:"headline" form:"headline"`
	Intro                         string
	Description                   string
	ProductLabel                  string `gorm:"column:productlabel"`
	ProductName                   string `gorm:"column:productname"`
	ProductDescription            string `gorm:"column:productdescription"`
	ProductAffiliateUrl           string `gorm:"column:productaffiliateurl"`
	Faq_Answer_1                  string
	Faq_Answer_2                  string
	Faq_Answer_3                  string
	Faq_Question_1                string
	Faq_Question_2                string
	Faq_Question_3                string
	HorizontalCardProductImageUrl string `gorm:"column:horizontalcardproductimageurl"`
	HorizontalCardProductImageAlt string `gorm:"column:horizontalcardproductimagealt"`
}
