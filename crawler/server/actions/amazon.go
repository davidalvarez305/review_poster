package actions

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/davidalvarez305/review_poster/crawler/server/types"
	"github.com/davidalvarez305/review_poster/crawler/server/utils"
)

type AmazonSearchResultsPage struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Link     string `json:"link"`
	Reviews  string `json:"reviews"`
	Price    string `json:"price"`
	Rating   string `json:"rating"`
	Category string `json:"category"`
}

type PAAAPI5Response struct {
	SearchResult types.SearchResult `json:"SearchResult"`
}

type AmazonPaapi5RequestBody struct {
	Marketplace string   `json:"Marketplace"`
	PartnerType string   `json:"PartnerType"`
	PartnerTag  string   `json:"PartnerTag"`
	Keywords    string   `json:"Keywords"`
	SearchIndex string   `json:"SearchIndex"`
	ItemCount   int      `json:"ItemCount"`
	Resources   []string `json:"Resources"`
}

func crawlPage(keyword, page string) ([]AmazonSearchResultsPage, error) {
	var crawledProducts []AmazonSearchResultsPage

	host := os.Getenv("P_HOST")
	username := os.Getenv("P_USERNAME")
	sessionId := fmt.Sprint(rand.Intn(1000000))
	path := username + sessionId + ":" + host

	u, err := url.Parse(path)

	if err != nil {
		fmt.Println("Error With URL Parse: ", err)
		return crawledProducts, err
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

	req, err := http.NewRequest("GET", page, nil)

	if err != nil {
		fmt.Println("Error With Proxy Request: ", err)
		return crawledProducts, err
	}

	req.Header.Set("Content-Type", "text/html")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while fetching Amazon SERP", err)
		return crawledProducts, err
	}
	defer resp.Body.Close()

	crawledProducts, err = parseHTML(resp.Body, keyword)

	if err != nil {
		fmt.Println("Error while parsing HTML.", err)
		return crawledProducts, err
	}

	return crawledProducts, nil
}

func ScrapeSearchResultsPage(keyword string) ([]AmazonSearchResultsPage, error) {
	var results []AmazonSearchResultsPage

	str := url.QueryEscape(keyword)

	wg := sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(page int) {
			serp := fmt.Sprintf("https://www.amazon.com/s?k=%s&s=review-rank&page=%v", str, page)

			products, err := crawlPage(keyword, serp)

			if err != nil {
				fmt.Printf("Error while crawling: %+v", err.Error())
			}

			results = append(results, products...)
			wg.Done()
		}(i)
	}

	wg.Wait()
	return results, nil
}

func SearchPaapi5Items(keyword string) (PAAAPI5Response, error) {
	var products PAAAPI5Response

	resources := []string{
		"Images.Primary.Medium",
		"ItemInfo.Title",
		"Offers.Listings.Price",
		"ItemInfo.ByLineInfo",
		"ItemInfo.Features",
		"ItemInfo.ProductInfo"}

	d := AmazonPaapi5RequestBody{
		Marketplace: "www.amazon.com",
		PartnerType: "Associates",
		PartnerTag:  os.Getenv("AMAZON_PARTNER_TAG"),
		Keywords:    keyword,
		SearchIndex: "All",
		ItemCount:   10,
		Resources:   resources,
	}

	body, err := json.Marshal(d)

	if err != nil {
		return products, err
	}

	method := "POST"
	service := "ProductAdvertisingAPI"
	url := "https://webservices.amazon.com/paapi5/searchitems"
	host := "webservices.amazon.com"
	region := os.Getenv("AWS_REGION")
	contentType := "application/json; charset=UTF-8"
	amazonTarget := "com.amazon.paapi5.v1.ProductAdvertisingAPIv1.SearchItems"
	contentEncoding := "amz-1.0"
	t := time.Now()
	amazonDate := utils.FormatShortDate(t)
	xAmazonDate := utils.FormatDate(t)
	canonicalUri := "/paapi5/searchitems"
	canonicalQuerystring := ""
	canonicalHeaders := utils.BuildCanonicalHeaders(contentType, contentEncoding, host, xAmazonDate, amazonTarget)
	credentialScope := amazonDate + "/" + region + "/" + service + "/" + "aws4_request"
	signedHeaders := "content-encoding;content-type;host;x-amz-date;x-amz-target"

	kSecret := os.Getenv("AWS_PAAPI_SECRET_KEY")
	kDate := utils.HMACSHA256([]byte("AWS4"+kSecret), []byte(amazonDate))
	kRegion := utils.HMACSHA256(kDate, []byte(region))
	kService := utils.HMACSHA256(kRegion, []byte(service))
	signingKey := utils.HMACSHA256(kService, []byte("aws4_request"))

	canonicalRequest := utils.BuildCanonicalString(method, canonicalUri, canonicalQuerystring, signedHeaders, canonicalHeaders, hex.EncodeToString(utils.MakeHash(sha256.New(), body)))
	stringToSign := utils.BuildStringToSign(xAmazonDate, credentialScope, canonicalRequest)
	signature, err := utils.BuildSignature(stringToSign, signingKey)
	if err != nil {
		fmt.Println("Error while building signature.")
		return products, err
	}

	authorizationHeader := "AWS4-HMAC-SHA256" + " Credential=" + os.Getenv("AWS_PAAPI_ACCESS_KEY") + "/" + credentialScope + " SignedHeaders=" + signedHeaders + " Signature=" + signature
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Request failed: ", err)
		return products, err
	}

	req.Header.Set("content-encoding", contentEncoding)
	req.Header.Set("content-type", contentType)
	req.Header.Set("host", host)
	req.Header.Set("x-amz-date", xAmazonDate)
	req.Header.Set("x-amz-target", amazonTarget)
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while fetching Amazon SERP", err)
		return products, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&products)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("STATUS CODE: %+v\n", resp.Status)
		return products, errors.New("request failed")
	}

	return products, nil
}

func parseHTML(r io.Reader, keyword string) ([]AmazonSearchResultsPage, error) {
	var results []AmazonSearchResultsPage
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Error trying to parse document.")
		return results, err
	}

	reviewsRegex := regexp.MustCompile("[0-9,]+")
	moneyRegex := regexp.MustCompile(`[\$]+?(\d+([,\.\d]+)?)`)
	amazonASIN := regexp.MustCompile(`(\/[A-Z0-9]{10,}\/)`)

	doc.Find(".sg-col-inner").Each(func(i int, s *goquery.Selection) {
		var product AmazonSearchResultsPage

		el, _ := s.Find("a").Attr("href")
		cond := amazonASIN.MatchString(el)

		if cond {
			name := strings.Join(strings.Split(strings.Split(el, "/")[1], "-"), " ")
			product.Name = name

			rating := strings.Split(s.Find(".a-icon-alt").Text(), " ")[0]
			product.Rating = rating

			link := strings.Split(el, "/")[3]
			product.Link = "https://amazon.com/dp/" + link + os.Getenv("AMAZON_TAG")

			image, _ := s.Find("img").Attr("src")
			product.Image = image

			product.Category = keyword

			if len(moneyRegex.FindAllString(s.Find(".a-size-base").Text(), 3)) > 0 {
				price := moneyRegex.FindAllString(s.Find(".a-size-base").Text(), 3)[0]
				product.Price = price
			}
			if len(reviewsRegex.FindAllString(s.Find(".a-size-base").Text(), 2)) > 0 {
				reviews := reviewsRegex.FindAllString(s.Find(".a-size-base").Text(), 3)[0]
				product.Reviews = reviews

			}
			results = append(results, product)
		}
	})

	return results, nil
}
