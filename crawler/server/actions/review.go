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

func CreateReviewPosts(categoryName, groupName string, dictionary types.DictionaryAPIResponse, sentences types.ContentAPIResponse) ([]models.ReviewPost, error) {
	var readyReviewPosts []models.ReviewPost

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{categoryName},
		},
	}

	googleKeywords, err := QueryGoogle(q)

	if err != nil {
		return readyReviewPosts, err
	}

	seedKeywords, err := GetSeedKeywords(googleKeywords)

	if err != nil {
		return readyReviewPosts, err
	}

	category, err := createOrFindCategory(categoryName, groupName)

	if err != nil {
		fmt.Printf("ERROR FINDING OR CREATING CATEGORY: %+v\n", err)
		return readyReviewPosts, err
	}

	subCategories, err := createSubCategories(seedKeywords, category)

	if err != nil {
		fmt.Printf("ERROR CREATING SUBCATEGORIES: %+v\n", err)
		return readyReviewPosts, err
	}

	wg := sync.WaitGroup{}
	for i := 0; i < len(seedKeywords)-1; i++ {
		wg.Add(1)
		go func(keywordNum int) {
			data, err := ScrapeSearchResultsPage(seedKeywords[keywordNum])

			if err != nil {
				fmt.Printf("ERROR SCRAPING: %+v\n", err)
				return
			}

			if len(data) == 0 {
				fmt.Println("Keyword: " + seedKeywords[keywordNum] + " - 0" + "\n")
				return
			}

			reviewPosts, err := createReviewPostsFactory(subCategories, seedKeywords[keywordNum], data, dictionary.Data, sentences.Data)

			if err != nil {
				fmt.Printf("ERROR INSERTING: %+v\n", err)
				return
			}

			readyReviewPosts = append(readyReviewPosts, reviewPosts...)

			total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", keywordNum+1, len(seedKeywords), seedKeywords[keywordNum], len(data))
			fmt.Println(total)

			fmt.Printf("Total Products = %v\n", len(reviewPosts))
			wg.Done()
		}(i)
	}

	wg.Wait()

	// Pull existing review posts so that I can exclude them from the slice of reviews to be created
	var existingReviewPosts []models.ReviewPost
	err = database.DB.Preload("Product").Find(&existingReviewPosts).Error
	if err != nil {
		fmt.Printf("ERROR FINDING EXISTING REVIEW POSTS: %+v\n", err)
		return readyReviewPosts, err
	}

	var reviewPostsTobeCreated []models.ReviewPost

	for _, post := range readyReviewPosts {
		exists := false

		// Ensure that no duplicates exist in DB
		for _, existingReviewPost := range existingReviewPosts {
			if post.Slug == existingReviewPost.Slug || post.Product.AffiliateUrl == existingReviewPost.ProductAffiliateUrl {
				exists = true
				break
			}
		}

		// If the first loop already found a duplicate, no need to execute the following code
		// Continue to next review post
		if exists {
			continue
		}

		// This code will make sure that not only the review post nor the product exist in the DB, but also don't exist in the current "TO BE CREATED" slice
		for _, acceptedReviewPost := range reviewPostsTobeCreated {
			if post.Slug == acceptedReviewPost.Slug || post.Product.AffiliateUrl == acceptedReviewPost.ProductAffiliateUrl {
				exists = true
				break
			}
		}

		if !exists {
			reviewPostsTobeCreated = append(reviewPostsTobeCreated, post)
		}
	}

	err = database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(&reviewPostsTobeCreated).Find(&reviewPostsTobeCreated).Error

	if err != nil {
		return reviewPostsTobeCreated, err
	}

	return reviewPostsTobeCreated, nil
}

func createReviewPostsFactory(subCategories []models.SubCategory, subCategoryName string, products []AmazonSearchResultsPage, dictionary []types.Word, sentences []types.Sentence) ([]models.ReviewPost, error) {
	var posts []models.ReviewPost

	for i := 0; i < len(products); i++ {
		slug := slug.Make(products[i].Name)
		replacedImage := strings.Replace(products[i].Image, "UL320", "UL640", 1)

		data := utils.GenerateContentUtil(products[i].Name, dictionary, sentences)

		var subCategoryId int
		for _, subcategory := range subCategories {
			if subcategory.Name == subCategoryName {
				subCategoryId = subcategory.ID
				break
			}
		}

		post := models.ReviewPost{
			Title:               data.ReviewPostTitle,
			SubCategoryID:       subCategoryId,
			Slug:                slug,
			Content:             data.ReviewPostContent,
			Headline:            data.ReviewPostHeadline,
			Intro:               data.ReviewPostIntro,
			Description:         data.ReviewPostDescription,
			Faq_Answer_1:        data.ReviewPostFaq_Answer_1,
			Faq_Answer_2:        data.ReviewPostFaq_Answer_2,
			Faq_Answer_3:        data.ReviewPostFaq_Answer_3,
			Faq_Question_1:      data.ReviewPostFaq_Question_1,
			Faq_Question_2:      data.ReviewPostFaq_Question_2,
			Faq_Question_3:      data.ReviewPostFaq_Question_3,
			ProductAffiliateUrl: products[i].Link,
			Product: &models.Product{
				AffiliateUrl:       products[i].Link,
				ProductPrice:       products[i].Price,
				ProductReviews:     products[i].Reviews,
				ProductRatings:     products[i].Rating,
				ProductImage:       replacedImage,
				ProductLabel:       data.ReviewPostProductLabel,
				ProductName:        products[i].Name,
				ProductDescription: data.ReviewPostProductDescription,
				ProductImageAlt:    strings.ToLower(products[i].Name),
			},
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// What this function does is it limits the requests to OpenAI to 1,000 every 60 seconds.
// The Requests Per Minute (RPM) Rate Limit is 3,500 but 1,000 requests is roughly $7.
/* func replaceContentWithChatGPT(posts []models.ReviewPost) []models.ReviewPost {
	var newReviewPosts []models.ReviewPost

	limit := 1000
	for i := 0; i < len(posts); {
		reviewPosts := asyncRequestsToOpenAI(posts[i:limit])
		newReviewPosts = append(newReviewPosts, reviewPosts...)
		i = limit
		limit += 1000

		fmt.Println("STARTING A 60-SECOND PAUSE")
		time.Sleep(60 * time.Second)
		fmt.Println("COMPLETED A 60-SECOND PAUSE")
	}

	return newReviewPosts
}

func asyncRequestsToOpenAI(posts []models.ReviewPost) []models.ReviewPost {
	var reviewPosts []models.ReviewPost
	wg := sync.WaitGroup{}
	for i := 0; i < len(posts)-1; i++ {
		wg.Add(1)
		go func(postNum int) {
			fmt.Println("REQUESTING CONTENT FROM OPEN AI")
			var post = posts[postNum]

			additionalContent, err := getAIGeneratedContent("What are people saying about the " + post.Product.ProductName)

			if err != nil {
				fmt.Printf("FAILED TO FETCH CONTENT FROM OPEN AI: %+v\n", err)
				return
			}

			FAQ_ONE, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + post.Faq_Answer_1)

			if err != nil {
				fmt.Printf("FAILED TO FETCH CONTENT FROM OPEN AI: %+v\n", err)
				return
			}

			FAQ_TWO, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + post.Faq_Answer_2)

			if err != nil {
				fmt.Printf("FAILED TO FETCH CONTENT FROM OPEN AI: %+v\n", err)
				return
			}

			FAQ_THREE, err := getAIGeneratedContent("Using college level writing, please re-write the following paragraph: " + post.Faq_Answer_3)

			if err != nil {
				fmt.Printf("FAILED TO FETCH CONTENT FROM OPEN AI: %+v\n", err)
				return
			}

			post.Content = post.Content + utils.GetAIResponse(additionalContent)
			post.Faq_Answer_1 = utils.GetAIResponse(FAQ_ONE)
			post.Faq_Answer_2 = utils.GetAIResponse(FAQ_TWO)
			post.Faq_Answer_3 = utils.GetAIResponse(FAQ_THREE)

			reviewPosts = append(reviewPosts, post)

			wg.Done()
		}(i)
	}

	wg.Wait()
	return reviewPosts
} */
