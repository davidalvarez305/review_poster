package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

type GoogleConfigData struct {
	Web struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthURI                 string   `json:"auth_uri"`
		TokenURI                string   `json:"token_uri"`
		AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
		JavascriptOrigins       []string `json:"javascript_origins"`
	} `json:"web"`
}

type TokenResponse struct {
	Access_Token string `json:"access_token"`
	Expires_In   string `json:"expires_in"`
	Scope        string `json:"scope"`
	Token_Type   string `json:"token_type"`
}

func (config *GoogleConfigData) GetGoogleCredentials() error {
	googlePath := os.Getenv("GOOGLE_JSON_PATH")
	path, err := ResolveServerPath()

	if err != nil {
		return err
	}

	jsonData, err := os.ReadFile(path + "/" + googlePath)

	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, &config)
}

func RequestGoogleAuthToken() error {
	config := &GoogleConfigData{}
	err := config.GetGoogleCredentials()

	if err != nil {
		return err
	}

	client := &http.Client{}

	url := config.Web.AuthURI
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("access_type", "offline")
	q.Add("approval_prompt", "force")
	q.Add("scope", "https://www.googleapis.com/auth/adwords")
	q.Add("client_id", config.Web.ClientID)
	q.Add("redirect_uri", config.Web.RedirectUris[0])
	q.Add("response_type", "code")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error while getting auth token", err)
		return err
	}
	defer resp.Body.Close()

	var data http.Response

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return err
	}

	return nil
}

func GetGoogleAccessToken(code string) (string, error) {
	config := &GoogleConfigData{}
	err := config.GetGoogleCredentials()

	if err != nil {
		return "", err
	}

	client := &http.Client{}

	url := config.Web.AuthURI
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", config.Web.ClientID)
	q.Add("client_secret", config.Web.ClientSecret)
	q.Add("redirect_uri", config.Web.RedirectUris[0])
	q.Add("scope", "https://www.googleapis.com/auth/adwords")
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	data := string(body)

	return data, nil

}

func RefreshAuthToken() (*oauth2.Token, error) {

	config := &GoogleConfigData{}
	err := config.GetGoogleCredentials()
	if err != nil {
		return nil, err
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	client := &http.Client{}

	url := config.Web.TokenURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("client_id", config.Web.ClientID)
	q.Add("client_secret", config.Web.ClientSecret)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}

	data := &oauth2.Token{}

	json.NewDecoder(resp.Body).Decode(data)

	return data, nil
}
