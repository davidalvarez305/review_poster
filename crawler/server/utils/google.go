package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/davidalvarez305/review_poster/crawler/server/types"
)

func GetGoogleCredentials() (types.GoogleConfigData, error) {
	data := types.GoogleConfigData{}

	path := os.Getenv("GOOGLE_JSON_PATH")

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not open Google JSON file.")
		return data, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Could not read Google JSON file.")
		return data, err
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error while trying to unmarshall JSON data.")
		return data, err
	}

	return data, nil
}
