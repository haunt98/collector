package confluence

import (
	"fmt"
	"regexp"
	"strings"
)

type Table struct {
	Headers []string
	Content [][]string
}

// | character conflict with | confluence table
// ignore | character inside [] confluence link
func normalizeVerticalCharacter(input string) string {
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

func ComposeTableFormat(table Table) string {
	var output string

	for _, header := range table.Headers {
		output += fmt.Sprintf("|| %s ", header)
	}
	output += "||\n"

	for _, row := range table.Content {
		for _, value := range row {
			normalizedValue := normalizeVerticalCharacter(value)
			output += fmt.Sprintf("| %s ", normalizedValue)
		}
		output += "|\n"
	}

	return output
}
