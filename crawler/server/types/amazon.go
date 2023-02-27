package types

type PAAAPI5Response struct {
	SearchResult SearchResult `json:"SearchResult"`
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
