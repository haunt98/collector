package confluence

import "fmt"

func FormatBold(input string) string {
	return fmt.Sprintf("*%s*", input)
}
