package actions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Prompt struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

func GetAdditionalContent(productName string) (OpenAIResponse, error) {
	var response OpenAIResponse
	url := "https://api.openai.com/v1/completions"

	prompt := Prompt{
		Model:       "text-davinci-003",
		Prompt:      "What are people saying about the " + productName,
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

	json.NewDecoder(resp.Body).Decode(&response)

	return response, nil
}
