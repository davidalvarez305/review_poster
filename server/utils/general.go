package utils

import "github.com/davidalvarez305/review_poster/server/types"

func GetAIResponse(response types.OpenAIResponse) string {
	var result string

	if len(response.Choices) > 0 {
		return response.Choices[0].Text
	}

	return result
}
