package actions

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/davidalvarez305/review_poster/crawler/server/utils"
	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"
)

func CreateNewReviewPost(input *AmazonSearchResultsPage, dictionary []types.Word, sentences []types.Sentence, subCategory models.SubCategory) (models.ReviewPost, error) {
	var post models.ReviewPost
	slug := slug.Make(input.Name)
	replacedImage := strings.Replace(input.Image, "UL320", "UL640", 1)

	additionalContent, err := GetAdditionalContent("What are people saying about the " + input.Name)

	if err != nil {
		return post, err
	}

	data := utils.GenerateContentUtil(input.Name, dictionary, sentences)

	FAQ_ONE, err := GetAdditionalContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_1)

	if err != nil {
		return post, err
	}

	FAQ_TWO, err := GetAdditionalContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_2)

	if err != nil {
		return post, err
	}

	FAQ_THREE, err := GetAdditionalContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_3)

	if err != nil {
		return post, err
	}

	// I'm creating this separately because there's a chance that the product already exists, in which case, it will be updated

	product := models.Product{
		AffiliateUrl:   input.Link,
		ProductPrice:   input.Price,
		ProductReviews: input.Reviews,
		ProductRatings: input.Rating,
		ProductImage:   input.Image,
	}

	err = database.DB.Clauses(clause.OnConflict{UpdateAll: true}).FirstOrCreate(&product).Error

	if err != nil {
		return post, err
	}

	post = models.ReviewPost{
		Title:               data.ReviewPostTitle,
		SubCategoryID:       subCategory.ID,
		ProductAffiliateUrl: product.AffiliateUrl,
		Slug:                slug,
		Content:             data.ReviewPostContent + additionalContent.Choices[0].Text,
		Headline:            data.ReviewPostHeadline,
		Intro:               data.ReviewPostIntro,
		Description:         data.ReviewPostDescription,
		ProductLabel:        data.ReviewPostProductLabel,
		ProductName:         input.Name,
		ProductDescription:  data.ReviewPostProductDescription,
		Faq_Answer_1:        FAQ_ONE.Choices[0].Text,
		Faq_Answer_2:        FAQ_TWO.Choices[0].Text,
		Faq_Answer_3:        FAQ_THREE.Choices[0].Text,
		Faq_Question_1:      data.ReviewPostFaq_Question_1,
		Faq_Question_2:      data.ReviewPostFaq_Question_2,
		Faq_Question_3:      data.ReviewPostFaq_Question_3,
		ProductImageUrl:     replacedImage,
		ProductImageAlt:     strings.ToLower(input.Name),
	}

	fmt.Printf("%+v\n", post)

	return post, nil
}

func InsertReviewPosts(groupName, categoryName, subCategoryName string, products AmazonSearchResultsPages, dictionary []types.Word, sentences []types.Sentence) error {
	var posts []models.ReviewPost

	subCategory, err := GetOrCreateSubCategory(categoryName, subCategoryName, groupName)

	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		p, err := CreateNewReviewPost(products[i], dictionary, sentences, *subCategory)

		if err != nil {
			continue
		}

		posts = append(posts, p)
	}

	err = database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(&posts).Error

	if err != nil {
		return err
	}

	return nil
}

func (products *AmazonSearchResultsPages) CreateReviewPosts(keyword, groupName string) error {
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
		data, err := ScrapeSearchResultsPage(seedKeywords[i])

		if err != nil {
			return err
		}

		if len(*data) == 0 {
			fmt.Println("Keyword: " + seedKeywords[i] + " - 0" + "\n")
			continue
		}

		err = InsertReviewPosts(groupName, keyword, seedKeywords[i], *data, dictionary.Data, sentences.Data)

		if err != nil {
			return err
		}

		results = append(results, *data...)

		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", i+1, len(seedKeywords), seedKeywords[i], len(*data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v\n", len(results))
	fmt.Println(productsTotal)

	products = &results
	return nil
}
