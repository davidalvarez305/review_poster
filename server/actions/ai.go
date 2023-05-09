package actions

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/davidalvarez305/review_poster/crawler/server/types"
)

func getAIGeneratedContent(promptMsg string) (types.OpenAIResponse, error) {
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
		os.Exit(1)
		return response, errors.New("request failed")
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response, nil
}
