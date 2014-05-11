package util

import (
		"strings"
		"regexp"
		// "fmt"
)

type Tokenizer int
const (
	TWhiteSpace Tokenizer = iota
	TLowerWhiteSpace = iota
)

var (
	whiteSpaceRe = regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")
)

func WhiteSpaceTokenize(doc string) []string {
	words := []string{}
	for _, word := range strings.Fields(doc) {
		words = append(words,  whiteSpaceRe.ReplaceAllString(word, ""))
	}
	return words
}

func LowerWhiteSpaceTokenize(doc string) []string {
	return WhiteSpaceTokenize(strings.ToLower(doc))
}



func Tokenize(t Tokenizer, doc string) []string {
	switch t {
	case TLowerWhiteSpace:
		return LowerWhiteSpaceTokenize(doc)
	}
	return WhiteSpaceTokenize(doc)

}