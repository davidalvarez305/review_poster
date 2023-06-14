package actions

import (
	"fmt"
	"sync"

	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/types"
	"gorm.io/gorm/clause"
)

func CreateReviewPosts(categoryName, groupName string, dictionary []models.Word, paragraphs []models.Paragraph) ([]models.ReviewPost, error) {
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

	// Unable to get any keywords -> stop
	if len(seedKeywords) == 0 {
		return readyReviewPosts, nil
	}

	// Generate more keywords from Open AI
	if len(seedKeywords) < 5 {
		seedKeywords = append(seedKeywords, GenerateKeywordsWithOpenAI(categoryName, seedKeywords)...)
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

	var wg sync.WaitGroup
	sem := make(chan struct{}, 20)

	for i := 0; i < len(seedKeywords); i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(keywordNum int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			data, err := ScrapeSearchResultsPage(seedKeywords[keywordNum])

			if err != nil {
				fmt.Printf("KEYWORD #%v - ERROR SCRAPING: %+v\n", keywordNum, err)
				return
			}

			if len(data) == 0 {
				fmt.Println("Keyword: " + seedKeywords[keywordNum] + " - 0" + "\n")
				return
			}

			reviewPosts, err := createReviewPostsFactory(subCategories, seedKeywords[keywordNum], data, dictionary, paragraphs)

			if err != nil {
				fmt.Printf("KEYWORD #%v - ERROR IN FACTORY: %+v\n", keywordNum, err)
				return
			}

			readyReviewPosts = append(readyReviewPosts, reviewPosts...)

			total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", keywordNum+1, len(seedKeywords), seedKeywords[keywordNum], len(data))
			fmt.Println(total)

			fmt.Printf("Total Products = %v\n", len(reviewPosts))
		}(i)
	}

	close(sem)
	wg.Wait()

	// Unable to crawl any review posts -> stop
	if len(readyReviewPosts) == 0 {
		return readyReviewPosts, nil
	}

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

	if len(reviewPostsTobeCreated) == 0 {
		return reviewPostsTobeCreated, nil
	}

	// Slicing the review posts because in the case of an insertion error, I don't want to lose all of the progress.
	var createdPosts []models.ReviewPost

	for i := 0; i < len(reviewPostsTobeCreated); i += 50 {
		end := i + 50
		if end > len(reviewPostsTobeCreated) {
			end = len(reviewPostsTobeCreated)
		}
		slicedList := reviewPostsTobeCreated[i:end]

		err = database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(&slicedList).Error

		if err != nil {
			continue
		}

		createdPosts = append(createdPosts, slicedList...)
	}

	return createdPosts, nil
}

func createReviewPostsFactory(subCategories []models.SubCategory, subCategoryName string, products []AmazonSearchResultsPage, dictionary []models.Word, paragraphs []models.Paragraph) ([]models.ReviewPost, error) {
	var posts []models.ReviewPost

	/* for i := 0; i < len(products); i++ {
		slug := slug.Make(products[i].Name)
		replacedImage := strings.Replace(products[i].Image, "UL320", "UL640", 1)

		data := utils.GenerateContentUtil(products[i].Name, dictionary, paragraphs)

		var subCategoryId int
		for _, subcategory := range subCategories {
			if subcategory.Name == subCategoryName {
				subCategoryId = subcategory.ID
				break
			}
		}

		if subCategoryId == 0 {
			break
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
	} */

	return posts, nil
}
