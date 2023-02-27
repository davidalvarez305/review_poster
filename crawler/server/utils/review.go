package utils

import (
	"regexp"
	"strings"
)

func CreateCategorySlug(str string) string {
	var final string
	r := regexp.MustCompile(`[a-z0-9]+`)
	res := r.FindAllString(str, -1)

	if len(res) > 0 {
		final = strings.Join(res, "-")
	}

	return final
}
