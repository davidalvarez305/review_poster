package actions

import (
	"fmt"

	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
)

func CreateReviewPosts(keyword, parent_group string) ([]types.AmazonSearchResultsPage, error) {
	var products []types.AmazonSearchResultsPage
	dictionary, err := PullContentDictionary()

	if err != nil {
		return products, err
	}

	sentences, err := PullDynamicContent()

	if err != nil {
		return products, err
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{keyword},
		},
	}

	KW_ARR, err := QueryGoogle(q)

	if err != nil {
		return products, err
	}

	seedKeywords, err := GetSeedKeywords(KW_ARR)

	if err != nil {
		return products, err
	}

	for i := 0; i < len(seedKeywords); i++ {
		data, err := ScrapeSearchResultsPage(seedKeywords[i])

		if err != nil {
			return products, err
		}

		if len(data) == 0 {
			fmt.Println("Keyword: " + seedKeywords[i] + "0" + "\n")
		}
		if len(data) > 0 {
			err := utils.InsertReviewPosts(parent_group, seedKeywords[i], data, dictionary, sentences)

			if err != nil {
				fmt.Printf("Error while trying to insert %s: %+v", seedKeywords[i], err)
				return products, err
			}

			products = append(products, data...)
		}
		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v\n", i+1, len(seedKeywords), seedKeywords[i], len(data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v", len(products))
	fmt.Println(productsTotal)

	return products, nil
}
