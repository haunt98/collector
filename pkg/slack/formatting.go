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

// https://api.slack.com/reference/surfaces/formatting#mentioning-users
const MentionChannel = "<!channel>"

func MentionUser(id string) string {
	return fmt.Sprintf("<@%s>", id)
}
