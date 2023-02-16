package actions

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
)

func RequestGoogleAuthToken() error {
	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return err
	}

	client := &http.Client{}

	url := config.Web.AuthURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
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
	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	client := &http.Client{}

	url := config.Web.AuthURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
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
		fmt.Println("Error while getting auth token", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body: ", err)
		return "", err
	}

	data := string(body)

	return data, nil

}

func RefreshAuthToken() (string, error) {

	type TokenResponse struct {
		Access_Token string `json:"access_token"`
		Expires_In   string `json:"expires_in"`
		Scope        string `json:"scope"`
		Token_Type   string `json:"token_type"`
	}

	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	client := &http.Client{}

	url := config.Web.TokenURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return "", err
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
		fmt.Println("Error while getting auth token: ", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Status)
	}

	var data TokenResponse

	json.NewDecoder(resp.Body).Decode(&data)

	return data.Access_Token, nil
}

func GetSeedKeywords(results types.GoogleKeywordResults) ([]string, error) {
	var data []string

	for i := 0; i < len(results.Results); i++ {
		compIndex, err := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.CompetitionIndex)
		if err != nil {
			return data, err
		}

		searchVol, err := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.AvgMonthlySearches)
		if err != nil {
			return data, err
		}

		keywordLength := len(strings.Split(results.Results[i].Text, " "))

		conditionOne := compIndex == 100
		conditionTwo := searchVol > 10000
		conditionThree := keywordLength >= 2 && keywordLength <= 4

		if conditionOne && conditionTwo && conditionThree {
			data = append(data, results.Results[i].Text)
		}
	}

	fmt.Println("Seed Keywords: ", len(data))

	return data, nil
}

func QueryGoogle(query types.GoogleQuery) (types.GoogleKeywordResults, error) {
	time.Sleep(1 * time.Second)
	var results types.GoogleKeywordResults

	authToken, err := RefreshAuthToken()

	if err != nil {
		fmt.Printf("Error refreshing token.")
		return results, err
	}

	googleCustomerID := os.Getenv("GOOGLE_CUSTOMER_ID")
	googleUrl := fmt.Sprintf("https://googleads.googleapis.com/v10/customers/%s:generateKeywordIdeas", googleCustomerID)
	developerToken := os.Getenv("GOOGLE_DEVELOPER_TOKEN")
	authorizationHeader := fmt.Sprintf("Bearer %s", authToken)

	client := &http.Client{}

	out, err := json.Marshal(query)

	if err != nil {
		return results, err
	}

	req, err := http.NewRequest("POST", googleUrl, bytes.NewBuffer(out))
	if err != nil {
		fmt.Println("Request failed: ", err)
		return results, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("developer-token", developerToken)
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying Google", err)
		return results, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&results)

	return results, nil
}

func GetCommercialKeywords(seedKeywords []string) ([]string, error) {
	var keywords []string
	for i := 0; i < len(seedKeywords); i++ {

		q := types.GoogleQuery{
			Pagesize: 1000,
			KeywordSeed: types.KeywordSeed{
				Keywords: [1]string{seedKeywords[i]},
			},
		}

		results, err := QueryGoogle(q)

		if err != nil {
			return keywords, err
		}

		k := utils.FilterCommercialKeywords(results, seedKeywords[i])
		keywords = append(keywords, k...)
	}

	fmt.Println("Commercial Keywords: ", len(keywords))

	return keywords, nil
}

func CrawlGoogleSERP(keywords string) ([]string, error) {
	var results []string

	str := strings.Join(strings.Split(keywords, " "), "+")

	serp := fmt.Sprintf("https://www.google.com/search?q=%s", str)

	host := os.Getenv("P_HOST")
	username := os.Getenv("P_USERNAME")
	sessionId := fmt.Sprint(rand.Intn(1000000))
	path := username + sessionId + ":" + host

	u, err := url.Parse(path)
	if err != nil {
		return results, err
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("GET", serp, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return results, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while fetching Google SERP", err)
		return results, err
	}
	defer resp.Body.Close()

	kws, err := utils.ParseGoogleSERP(resp.Body)
	results = kws

	if err != nil {
		fmt.Println("Error while parsing HTML.")
		return kws, err
	}

	return results, nil
}
