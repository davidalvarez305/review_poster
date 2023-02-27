package models

type Group struct {
	Base
	Categories []*Category `json:"categories" form:"categories"`
}
