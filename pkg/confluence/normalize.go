package confluence

import (
	"regexp"
	"strings"
)

// | character conflict with | confluence table
// ignore | character inside [] confluence link
func NormalizeVerticalCharacter(input string) string {
	output := strings.ReplaceAll(input, "|", `\|`)

	// change back | inside []
	regex := regexp.MustCompile(`\[.*?\\\|.*?]`)
	if regex.MatchString(output) {
		links := regex.FindAllString(output, -1)

		for _, link := range links {
			newLink := strings.ReplaceAll(link, `\|`, "|")
			output = strings.ReplaceAll(output, link, newLink)
		}
	}

	return output
}

// * is confluence bullet list

func NormalizeList(input string) string {
	input = strings.ReplaceAll(input, "• ", "* ")
	input = strings.ReplaceAll(input, "- ", "* ")

	return input
}

func TitleList(input string) string {
	regex := regexp.MustCompile(`\* .`)
	if !regex.MatchString(input) {
		return input
	}

	replaceFn := func(s string) string {
		return strings.ToUpper(s)
	}

	input = regex.ReplaceAllStringFunc(input, replaceFn)

	return input
}
