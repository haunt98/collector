package slack

import (
	"fmt"
	"strings"
)

// https://api.slack.com/reference/surfaces/formatting#mentioning-users
func NormalizeUser(input string, users []User) string {
	for _, user := range users {
		input = strings.ReplaceAll(input, fmt.Sprintf("<@%s>",
			user.ID), fmt.Sprintf("@%s", user.Profile.DisplayName))
	}

	return input
}
