package actions

import (
	"fmt"
	"strings"
	"sync"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/davidalvarez305/review_poster/crawler/server/utils"
	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"
)

func CreateReviewPosts(categoryName, groupName string, dictionary types.DictionaryAPIResponse, sentences types.ContentAPIResponse) ([]AmazonSearchResultsPage, error) {
	var products []AmazonSearchResultsPage

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{categoryName},
		},
	}

	googleKeywords, err := QueryGoogle(q)

	if err != nil {
		return products, err
	}

	seedKeywords, err := GetSeedKeywords(googleKeywords)

	if err != nil {
		return products, err
	}

	category, err := createOrFindCategory(categoryName, groupName)

	if err != nil {
		return products, err
	}

	subCategories, err := createSubCategories(seedKeywords, category)

	if err != nil {
		fmt.Printf("ERROR CREATING SUBCATEGORIES: %+v\n", err)
		return products, err
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(seedKeywords)-1; i++ {
		wg.Add(1)
		go func(keywordNum int) {
			data, err := ScrapeSearchResultsPage(seedKeywords[keywordNum])

			if err != nil {
				fmt.Printf("ERROR SCRAPING: %+v\n", err)
			}

			if len(data) == 0 {
				fmt.Println("Keyword: " + seedKeywords[keywordNum] + " - 0" + "\n")
			}

			err = insertReviewPosts(subCategories, seedKeywords[keywordNum], data, dictionary.Data, sentences.Data)

			if err != nil {
				fmt.Printf("ERROR INSERTING: %+v\n", err)
			}

			products = append(products, data...)

			total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", keywordNum+1, len(seedKeywords), seedKeywords[keywordNum], len(data))
			fmt.Println(total)

			fmt.Printf("Total Products = %v\n", len(products))
			wg.Done()
		}(i)
	}

	wg.Wait()
	return products, nil
}

func insertReviewPosts(subCategories []models.SubCategory, subCategoryName string, products []AmazonSearchResultsPage, dictionary []types.Word, sentences []types.Sentence) error {
	var posts []models.ReviewPost

	wg := sync.WaitGroup{}
	for i := 0; i < len(products)-1; i++ {
		wg.Add(1)
		go func(productNum int) {
			p, err := assembleReviewPost(products[productNum], dictionary, sentences, subCategories, subCategoryName)

			if err != nil {
				fmt.Printf("ERROR CREATING NEW REVIEW POST: %+v\n", err)
			}

			err = database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(&p).Error

			if err != nil {
				fmt.Printf("ERROR SAVING NEW REVIEW POST: %+v\n", err)
			}

			fmt.Printf("Product successfully crawled: %+v\n", p.Title)
			posts = append(posts, p)
			wg.Done()
		}(i)
	}

	wg.Wait()

	return nil
}

func assembleReviewPost(input AmazonSearchResultsPage, dictionary []types.Word, sentences []types.Sentence, subCategories []models.SubCategory, keyword string) (models.ReviewPost, error) {
	var post models.ReviewPost
	slug := slug.Make(input.Name)
	replacedImage := strings.Replace(input.Image, "UL320", "UL640", 1)

	additionalContent, err := getAIGeneratedContent("What are people saying about the " + input.Name)

	if err != nil {
		return post, err
	}

	data := utils.GenerateContentUtil(input.Name, dictionary, sentences)

	FAQ_ONE, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_1)

	if err != nil {
		return post, err
	}

	FAQ_TWO, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_2)

	if err != nil {
		return post, err
	}

	FAQ_THREE, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + data.ReviewPostFaq_Answer_3)

	if err != nil {
		return post, err
	}

	// I'm creating this separately because there's a chance that the product already exists, in which case, it will be updated

	product := models.Product{
		AffiliateUrl:       input.Link,
		ProductPrice:       input.Price,
		ProductReviews:     input.Reviews,
		ProductRatings:     input.Rating,
		ProductImage:       replacedImage,
		ProductLabel:       data.ReviewPostProductLabel,
		ProductName:        input.Name,
		ProductDescription: data.ReviewPostProductDescription,
		ProductImageAlt:    strings.ToLower(input.Name),
	}

	err = database.DB.Clauses(clause.OnConflict{UpdateAll: true}).Save(&product).Error

	if err != nil {
		return post, err
	}

	var subCategoryId int
	for _, subcategory := range subCategories {
		if subcategory.Name == keyword {
			subCategoryId = subcategory.ID
			break
		}
	}

	post = models.ReviewPost{
		Title:               data.ReviewPostTitle,
		SubCategoryID:       subCategoryId,
		ProductAffiliateUrl: input.Link,
		Slug:                slug,
		Content:             data.ReviewPostContent + utils.GetAIResponse(additionalContent),
		Headline:            data.ReviewPostHeadline,
		Intro:               data.ReviewPostIntro,
		Description:         data.ReviewPostDescription,
		Faq_Answer_1:        utils.GetAIResponse(FAQ_ONE),
		Faq_Answer_2:        utils.GetAIResponse(FAQ_TWO),
		Faq_Answer_3:        utils.GetAIResponse(FAQ_THREE),
		Faq_Question_1:      data.ReviewPostFaq_Question_1,
		Faq_Question_2:      data.ReviewPostFaq_Question_2,
		Faq_Question_3:      data.ReviewPostFaq_Question_3,
	}

	return post, nil
}
