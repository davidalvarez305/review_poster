package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/davidalvarez305/review_poster/server/types"
)

func GetGoogleCredentials() (types.GoogleConfigData, error) {
	data := types.GoogleConfigData{}

	path, err := filepath.Abs(os.Getenv("GOOGLE_JSON_PATH"))

	if err != nil {
		return data, err
	}

	file, err := os.Open(path)

	if err != nil {
		return data, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		return data, err
	}

	return data, nil
}
