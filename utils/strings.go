package utils

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var multiWhitespace = regexp.MustCompile(`\s\p{Zs}+`)

// StringOperation Type for functions that operate on strings and return a string
type StringOperation func(string) string

// GenerateReplacePrefix creates a function for replacing a prefix with another prefix
func GenerateReplacePrefix(oldPrefix string, newPrefix string) StringOperation {
	return func(original string) string {
		if strings.HasPrefix(original, oldPrefix) {
			return newPrefix + original[len(oldPrefix):]
		}
		return original
	}
}

// GenerateMockReplacePrefix creates no-op function that takes a string and returns the same string
func GenerateMockReplacePrefix() StringOperation {
	return func(original string) string {
		return original
	}
}

func StripWhitespace(inputString string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, inputString)
}

func NormalizeString(inputString string) (string, int, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	return transform.String(t, inputString)
}

func CompressWhitespace(input string) string {
	return multiWhitespace.ReplaceAllString(input, " ")
}
