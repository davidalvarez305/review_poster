package utils

import (
	"math/rand"
	"regexp"
	"strings"

	"github.com/davidalvarez305/review_poster/server/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func filterSentences(sentences []types.Sentence, paragraph string) []string {
	var s []string
	for i := 0; i < len(sentences); i++ {
		if sentences[i].Paragraph.Name == paragraph {

			s = append(s, sentences[i].Sentence)
		}
	}
	return s
}

func processSentence(productName, sentence string, dictionary []types.Word) string {
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

func switchWords(matchedWord string, dictionary []types.Word) string {
	for i := 0; i < len(dictionary); i++ {
		if dictionary[i].Tag == matchedWord {
			matchedWord = dictionary[i].Synonyms[rand.Intn(len(dictionary[i].Synonyms))].Synonym
		}
	}
	return matchedWord
}

func spinnerFunction(productName, matchedWord string, dictionary []types.Word) string {
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

func selectRandomSentences(productName string, sentences []types.ProcessedContent, dictionary []types.Word) types.FinalizedContent {
	var content types.FinalizedContent
	for i := 0; i < len(sentences); i++ {
		content.ReviewPostTitle = processSentence(productName, sentences[i].ReviewPostTitle[rand.Intn(len(sentences[i].ReviewPostTitle))], dictionary)
		content.ReviewPostContent = processSentence(productName, sentences[i].ReviewPostContent[rand.Intn(len(sentences[i].ReviewPostContent))], dictionary)
		content.ReviewPostHeadline = processSentence(productName, sentences[i].ReviewPostHeadline[rand.Intn(len(sentences[i].ReviewPostHeadline))], dictionary)
		content.ReviewPostIntro = processSentence(productName, sentences[i].ReviewPostIntro[rand.Intn(len(sentences[i].ReviewPostIntro))], dictionary)
		content.ReviewPostDescription = processSentence(productName, sentences[i].ReviewPostDescription[rand.Intn(len(sentences[i].ReviewPostDescription))], dictionary)
		content.ReviewPostProductLabel = processSentence(productName, sentences[i].ReviewPostProductLabel[rand.Intn(len(sentences[i].ReviewPostProductLabel))], dictionary)
		content.ReviewPostProductDescription = processSentence(productName, sentences[i].ReviewPostProductDescription[rand.Intn(len(sentences[i].ReviewPostProductDescription))], dictionary)
		content.ReviewPostFaq_Answer_1 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_1[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_1))], dictionary)
		content.ReviewPostFaq_Answer_2 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_2[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_2))], dictionary)
		content.ReviewPostFaq_Answer_3 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_3[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_3))], dictionary)
		content.ReviewPostFaq_Question_1 = processSentence(productName, sentences[i].ReviewPostFaq_Question_1[rand.Intn(len(sentences[i].ReviewPostFaq_Question_1))], dictionary)
		content.ReviewPostFaq_Question_2 = processSentence(productName, sentences[i].ReviewPostFaq_Question_2[rand.Intn(len(sentences[i].ReviewPostFaq_Question_2))], dictionary)
		content.ReviewPostFaq_Question_3 = processSentence(productName, sentences[i].ReviewPostFaq_Question_3[rand.Intn(len(sentences[i].ReviewPostFaq_Question_3))], dictionary)
	}
	return content
}

// In the future, this might have to be re-worked.
func GenerateContentUtil(productName string, dictionary []types.Word, sentences []types.Sentence) types.FinalizedContent {
	var content []types.ProcessedContent
	var finalContent types.FinalizedContent

	for i := 0; i < len(sentences); i++ {
		a := filterSentences(sentences, "ReviewPostTitle")
		b := filterSentences(sentences, "ReviewPostContent")
		c := filterSentences(sentences, "ReviewPostHeadline")
		d := filterSentences(sentences, "ReviewPostIntro")
		e := filterSentences(sentences, "ReviewPostDescription")
		f := filterSentences(sentences, "ReviewPostProductLabel")
		g := filterSentences(sentences, "ReviewPostProductDescription")
		h := filterSentences(sentences, "ReviewPostFaq_Answer_1")
		j := filterSentences(sentences, "ReviewPostFaq_Answer_2")
		k := filterSentences(sentences, "ReviewPostFaq_Answer_3")
		l := filterSentences(sentences, "ReviewPostFaq_Question_1")
		m := filterSentences(sentences, "ReviewPostFaq_Question_2")
		n := filterSentences(sentences, "ReviewPostFaq_Question_3")
		var final = types.ProcessedContent{
			ReviewPostTitle:              a,
			ReviewPostContent:            b,
			ReviewPostHeadline:           c,
			ReviewPostIntro:              d,
			ReviewPostDescription:        e,
			ReviewPostProductLabel:       f,
			ReviewPostProductDescription: g,
			ReviewPostFaq_Answer_1:       h,
			ReviewPostFaq_Answer_2:       j,
			ReviewPostFaq_Answer_3:       k,
			ReviewPostFaq_Question_1:     l,
			ReviewPostFaq_Question_2:     m,
			ReviewPostFaq_Question_3:     n,
		}
		content = append(content, final)
	}
	finalContent = selectRandomSentences(productName, content, dictionary)
	return finalContent
}
