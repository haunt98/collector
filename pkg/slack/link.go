package slack

import (
	"regexp"
	"strings"
)

type Link struct {
	Original    string
	URL         string
	Description string
}

// https://api.slack.com/reference/surfaces/formatting#linking-urls
func ExtractLinks(input string) ([]Link, bool) {
	// <url|description>
	regex := regexp.MustCompile(`<((?:http|www\.)\S+?)\|(.+?)>`)
	if regex.MatchString(input) {
		submatches := regex.FindAllStringSubmatch(input, -1)

		links := make([]Link, 0, len(submatches))
		for i := range submatches {
			link := Link{
				Original:    submatches[i][0],
				URL:         submatches[i][1],
				Description: submatches[i][2],
			}

			links = append(links, link)
		}

		return links, true
	}

	// <url>
	regex = regexp.MustCompile(`<((?:http|www\.)[^\s|]+?)>`)
	if regex.MatchString(input) {
		submatches := regex.FindAllStringSubmatch(input, -1)

		links := make([]Link, 0, len(submatches))
		for i := range submatches {
			link := Link{
				Original: submatches[i][0],
				URL:      submatches[i][1],
			}

			links = append(links, link)
		}

		return links, true
	}

	// url
	regex = regexp.MustCompile(`((?:http|www\.)[^\s|]+?)(?:\s|$)`)
	if regex.MatchString(input) {
		submatches := regex.FindAllStringSubmatch(input, -1)

		links := make([]Link, 0, len(submatches))
		for i := range submatches {
			link := Link{
				Original: strings.TrimSpace(submatches[i][0]),
				URL:      submatches[i][1],
			}

			links = append(links, link)
		}

		return links, true
	}

	return nil, false
}
