package slack

import "strings"

// https://api.slack.com/reference/surfaces/formatting#visual-styles

func RemoveBold(input string) string {
	return strings.ReplaceAll(input, "*", "")
}
