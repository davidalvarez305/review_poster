package actions

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/davidalvarez305/review_poster/crawler/server/utils"
	"gorm.io/gorm/clause"
)

func CreateNewReviewPost(input *AmazonSearchResultsPage, dictionary []types.Dictionary, sentences []types.DynamicContent, subCategory models.SubCategory) models.ReviewPost {
	name := strings.Join(strings.Split(strings.ToLower(input.Name), " "), "-")
	slug := utils.CreateCategorySlug(name)
	replacedImage := strings.Replace(input.Image, "UL320", "UL640", 1)

	data := utils.GenerateContentUtil(input.Name, dictionary, sentences)

	post := models.ReviewPost{
		Title:              data.ReviewPostTitle,
		SubCategory:        &subCategory,
		Slug:               slug,
		Content:            data.ReviewPostContent,
		Headline:           data.ReviewPostHeadline,
		Intro:              data.ReviewPostIntro,
		Description:        data.ReviewPostDescription,
		ProductLabel:       data.ReviewPostProductLabel,
		ProductName:        input.Name,
		ProductDescription: data.ReviewPostProductDescription,
		Product: &models.Product{
			AffiliateUrl:   input.Link,
			ProductPrice:   input.Price,
			ProductReviews: input.Reviews,
			ProductRatings: input.Rating,
			ProductImage:   input.Image,
		},
		Faq_Answer_1:    data.ReviewPostFaq_Answer_1,
		Faq_Answer_2:    data.ReviewPostFaq_Answer_2,
		Faq_Answer_3:    data.ReviewPostFaq_Answer_3,
		Faq_Question_1:  data.ReviewPostFaq_Question_1,
		Faq_Question_2:  data.ReviewPostFaq_Question_2,
		Faq_Question_3:  data.ReviewPostFaq_Question_3,
		ProductImageUrl: replacedImage,
		ProductImageAlt: strings.ToLower(input.Name),
	}

	return post
}

func InsertReviewPosts(groupName, categoryName, subCategoryName string, products AmazonSearchResultsPages, dictionary []types.Dictionary, sentences []types.DynamicContent) error {
	var posts []models.ReviewPost

	subCategory := &SubCategory{}

	err := subCategory.GetOrCreateSubCategory(categoryName, subCategoryName, groupName)

	if err != nil {
		return err
	}

	for i := 0; i < len(products); i++ {
		p := CreateNewReviewPost(products[i], dictionary, sentences, *subCategory.SubCategory)
		posts = append(posts, p)
	}

	db := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(posts)

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (products *AmazonSearchResultsPages) CreateReviewPosts(keyword, groupName, categoryName string) error {
	dictionary, err := PullContentDictionary()
	googleKeywords := &GoogleKeywordResults{}
	var results AmazonSearchResultsPages

	if err != nil {
		return err
	}

	sentences, err := PullDynamicContent()

	if err != nil {
		return err
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{keyword},
		},
	}

	err = googleKeywords.QueryGoogle(q)

	if err != nil {
		return err
	}

	seedKeywords, err := GetSeedKeywords(googleKeywords)

	if err != nil {
		return err
	}

	for i := 0; i < len(seedKeywords); i++ {
		var data AmazonSearchResultsPages
		err := data.ScrapeSearchResultsPage(seedKeywords[i])

		if err != nil {
			return err
		}

		if len(data) == 0 {
			fmt.Println("Keyword: " + seedKeywords[i] + "0" + "\n")
		}
		if len(data) > 0 {
			err := InsertReviewPosts(groupName, categoryName, seedKeywords[i], data, dictionary, sentences)

			if err != nil {
				fmt.Printf("Error while trying to insert %s: %+v", seedKeywords[i], err)
				return err
			}

			results = append(results, data...)
		}
		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", i+1, len(seedKeywords), seedKeywords[i], len(data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v", len(results))
	fmt.Println(productsTotal)

	products = &results
	return nil
}
