package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/davidalvarez305/review_poster/crawler/server/types"
)

func PullDynamicContent() (types.ContentAPIResponse, error) {
	var content types.ContentAPIResponse
	contentApi := os.Getenv("DYNAMIC_CONTENT_API") + "content?template=ReviewPost"

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentApi, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return content, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("AUTH_HEADER_STRING"))
	req.Header.Set("X-Secret-Agent", os.Getenv("X_SECRET_AGENT"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying the dynamic content endpoint.", err)
		return content, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		return content, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&content)
	return content, nil
}

func PullContentDictionary() (types.DictionaryAPIResponse, error) {
	var content types.DictionaryAPIResponse
	contentApi := os.Getenv("DYNAMIC_CONTENT_API") + "dictionary"

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentApi, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return content, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("AUTH_HEADER_STRING"))
	req.Header.Set("X-Secret-Agent", os.Getenv("X_SECRET_AGENT"))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying the dictionary endpoint.", err)
		return content, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		return content, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&content)
	return content, nil
}