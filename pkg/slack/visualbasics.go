package slack

import (
	"fmt"
	"strings"
)

// https://api.slack.com/reference/surfaces/formatting#visual-styles

func RemoveBold(input string) string {
	return strings.ReplaceAll(input, "*", "")
}

func AddBold(input string) string {
	return fmt.Sprintf("*%s*", input)
}
