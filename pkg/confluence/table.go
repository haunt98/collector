package confluence

import "fmt"

type Table struct {
	Headers []string
	Content [][]string
}

func ComposeTableFormat(table Table) string {
	var output string

	for _, header := range table.Headers {
		output += fmt.Sprintf("|| %s ", header)
	}
	output += "||\n"

	for _, row := range table.Content {
		for _, value := range row {
			output += fmt.Sprintf("| %s ", value)
		}
		output += "|\n"
	}

	return output
}
