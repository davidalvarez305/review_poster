package types

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

type KeywordSeed struct {
	Keywords [1]string `json:"keywords"`
}

type GoogleQuery struct {
	Pagesize    int         `json:"pageSize"`
	KeywordSeed KeywordSeed `json:"keywordSeed"`
}

type MonthlySearchVolume struct {
	Month           string `json:"month"`
	Year            string `json:"year"`
	MonthlySearches string `json:"monthlySearches"`
}

type keywordIdeaMetrics struct {
	Competition            string                `json:"competition"`
	MonthlySearchVolume    []MonthlySearchVolume `json:"monthlySearchVolumes"`
	AvgMonthlySearches     string                `json:"avgMonthlySearches"`
	CompetitionIndex       string                `json:"competitionIndex"`
	LowTopOfPageBidMicros  string                `json:"lowTopOfPageBidMicros"`
	HighTopOfPageBidMicros string                `json:"highTopOfPageBidMicros"`
}

type GoogleResult struct {
	KeywordIdeaMetrics keywordIdeaMetrics `json:"keywordIdeaMetrics"`
	Text               string             `json:"text"`
}

type GoogleKeywordResults struct {
	Results []GoogleResult `json:"results"`
}
