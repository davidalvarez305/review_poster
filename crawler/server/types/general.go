package types

type CreatePostInput struct {
	ProductName    string `json:"productName"`
	Category       string `json:"category"`
	ProductURL     string `json:"productUrl"`
	ImageURL       string `json:"imageUrl"`
	ProductReviews string `json:"productReviews"`
	ProductPrice   string `json:"productPrice"`
	ProductRating  string `json:"productRating"`
}

type Keyword struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type CreateReviewPostsInput struct {
	Keyword       string `json:"keyword" form:"keyword"`
	GroupName     string `json:"groupName" form:"groupName"`
	CatergoryName string `json:"categoryName" form:"categoryName"`
}
