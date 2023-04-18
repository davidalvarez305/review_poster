package models

type ReviewPost struct {
	ID                  int          `json:"id" form:"id"`
	Title               string       `json:"title" form:"title"`
	Slug                string       `gorm:"unique" json:"slug" form:"slug"`
	Content             string       `json:"content" form:"content"`
	SubCategoryID       int          `json:"sub_category_id" form:"sub_category_id"`
	SubCategory         *SubCategory `gorm:"not null;column:sub_category_id;foreignKey:SubCategoryID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"sub_category" form:"sub_category"`
	Headline            string       `json:"headline" form:"headline"`
	Intro               string       `json:"intro" form:"intro"`
	Description         string       `json:"description" form:"description"`
	ProductAffiliateUrl string       `gorm:"column:product_affiliate_url" json:"product_affiliate_url" form:"product_affiliate_url"`
	Product             *Product     `gorm:"not null;foreignKey:ProductAffiliateUrl;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"product" form:"product"`
	Faq_Answer_1        string       `gorm:"column:faq_answer_1" json:"faq_answer_1" form:"faq_answer_1"`
	Faq_Answer_2        string       `gorm:"column:faq_answer_2" json:"faq_answer_2" form:"faq_answer_2"`
	Faq_Answer_3        string       `gorm:"column:faq_answer_3" json:"faq_answer_3" form:"faq_answer_3"`
	Faq_Question_1      string       `gorm:"column:faq_question_1" json:"faq_question_1" form:"faq_question_1"`
	Faq_Question_2      string       `gorm:"column:faq_question_2" json:"faq_question_2" form:"faq_question_2"`
	Faq_Question_3      string       `gorm:"column:faq_question_3" json:"faq_question_3" form:"faq_question_3"`
}
