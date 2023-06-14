package utils

import (
	"math/rand"
	"regexp"
	"strings"

	"github.com/davidalvarez305/review_poster/server/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func processSentence(productName, sentence string, dictionary []models.Word) string {
	var s string

	r := regexp.MustCompile(`(\([#@]\w+:[A-Z]+)\)|(\([#@]\w+)\)`)

	m := r.FindAllString(sentence, -1)

	for i := 0; i < len(m); i++ {
		if i == 0 {
			s = sentence
		}
		switched := strings.Replace(s, m[i], spinnerFunction(productName, m[i], dictionary), -1)
		s = switched
	}
	return s
}

func switchWords(matchedWord string, dictionary []models.Word) string {
	for i := 0; i < len(dictionary); i++ {
		if dictionary[i].Tag == matchedWord {
			matchedWord = dictionary[i].Synonyms[rand.Intn(len(dictionary[i].Synonyms))].Synonym
		}
	}
	return matchedWord
}

func spinnerFunction(productName, matchedWord string, dictionary []models.Word) string {
	makeTitle := cases.Title(language.English)
	if matchedWord == "(@ProductName)" {
		matchedWord = productName
		return matchedWord
	} else {
		splitStr := strings.Split(matchedWord, ":")
		if len(splitStr) == 2 {
			s := splitStr[0] + ")"
			matchedWord = switchWords(s, dictionary)
			if splitStr[1] == "UU)" {
				matchedWord = makeTitle.String(matchedWord)
			}
			if splitStr[1] == "U)" {
				ss := strings.Split(matchedWord, "")
				ss[0] = strings.ToUpper(ss[0])
				matchedWord = strings.Join(ss, "")
			}
		} else {
			matchedWord = switchWords(matchedWord, dictionary)
		}
	}
	return matchedWord
}

func selectRandomSentences(productName string, paragraphs []models.Paragraph, dictionary []models.Word) map[string]string {
	content := make(map[string]string, len(paragraphs))
	for i := 0; i < len(paragraphs); i++ {
		if len(paragraphs[i].Sentences) == 0 {
			continue
		}
		randomSentence := paragraphs[i].Sentences[rand.Intn(len(paragraphs[i].Sentences))].Sentence
		content[paragraphs[i].Name] = randomSentence
	}
	return content
}

// In the future, this might have to be re-worked.
func GenerateContentUtil(productName string, dictionary []models.Word, paragraphs []models.Paragraph) map[string]string {
	return selectRandomSentences(productName, paragraphs, dictionary)
}
