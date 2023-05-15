package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/davidalvarez305/review_poster/server/types"
)

func QueryOpenAI(promptMsg string) (types.OpenAIResponse, error) {
	var response types.OpenAIResponse
	url := "https://api.openai.com/v1/completions"

	prompt := types.Prompt{
		Model:       "text-davinci-003",
		Prompt:      promptMsg,
		Temperature: 0.6,
		MaxTokens:   2000,
	}

	data, err := json.Marshal(prompt)

	if err != nil {
		return response, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

	if err != nil {
		return response, err
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + os.Getenv("OPEN_AI_KEY")},
	}

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		fmt.Printf("ERR BODY: %+v\n", resp.Body)
		return response, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response, nil
}

func GenerateKeywordsWithOpenAI(categoryName string, seedKeywords []string) []string {
	var generatedKeywords []string
	var wg sync.WaitGroup

	for _, seedKeyword := range seedKeywords {
		wg.Add(1)
		go func(keyword string) {
			defer func() {
				wg.Done()
			}()
			response, err := QueryOpenAI(fmt.Sprintf("Please give me a list of the top brands on Amazon.com for the %s category. Please do not enumarate the list, and list all entries lowercased, and separated by lines.", keyword))

			if err != nil {
				fmt.Printf("ERROR GETTING KEYWORDS FROM OPEN AI: %+v\n", err)
				return
			}

			if len(response.Choices) > 0 {
				var filteredEntries []string

				for _, result := range strings.Split(response.Choices[0].Text, "\n") {
					if len(result) > 0 {
						filteredEntries = append(filteredEntries, result+" "+categoryName)
					}
				}

				generatedKeywords = append(generatedKeywords, filteredEntries...)
			}
		}(seedKeyword)
	}

	wg.Wait()

	return generatedKeywords
}
