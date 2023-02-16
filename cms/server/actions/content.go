package actions

import (
	"fmt"

	"github.com/davidalvarez305/content_go/server/database"
)

type GetDictionary struct {
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Synonyms string `json:"synonyms"`
}

type Dictionary []*GetDictionary

type GetContent struct {
	Sentences string `json:"sentences"`
	Template  string `json:"template"`
	Paragraph string `json:"paragraph"`
	Order     int    `json:"order"`
}

type Content []*GetContent

func (content *Content) GetContent(template string) error {
	query := fmt.Sprintf(
		`SELECT t.name AS template, p.name AS paragraph, p.order AS order, string_agg(s.sentence, '///') AS sentences
		FROM sentence AS s
		INNER JOIN template AS t
		ON (s.template_id = t.id)
		INNER JOIN paragraph AS p
		ON (s.paragraph_id = p.id)
		WHERE '%s' = t.name
		GROUP BY t.name, p.name, p.order;`, template)

	result := database.DB.Raw(query).Scan(&content)

	return result.Error
}

func (dictionary *Dictionary) GetDictionary() error {
	query :=
		`SELECT w.name, w.tag, string_agg(d.synonym, '///') AS synonyms
		FROM word AS w
		INNER JOIN synonym AS d
		ON (w.id = d.word_id)
		GROUP BY w.name, w.tag;`

	result := database.DB.Raw(query).Scan(&dictionary)

	return result.Error
}
