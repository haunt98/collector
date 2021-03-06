package confluence

import (
	"regexp"
	"strings"
)

// * is confluence bullet list

func normalizeList(input string) string {
	// Bullet in start of line
	regex := regexp.MustCompile(`^[•\-]\s+`)
	if regex.MatchString(input) {
		input = regex.ReplaceAllString(input, "* ")
	}

	// Bullet before newline
	regex = regexp.MustCompile(`\n\s*?[•\-]\s+`)
	if regex.MatchString(input) {
		input = regex.ReplaceAllString(input, "\n* ")
	}

	return input
}

func titleList(input string) string {
	regex := regexp.MustCompile(`\* .`)
	if !regex.MatchString(input) {
		return input
	}

	input = regex.ReplaceAllStringFunc(input, strings.ToUpper)

	return input
}

func ComposeListFormat(input string) string {
	input = normalizeList(input)
	input = titleList(input)

	return input
}
