package scrum

import (
	"collector/pkg/confluence"
	"collector/pkg/slack"
	"regexp"
	"strings"
)

const (
	domainTitle   = "Domain"
	beforeTitle   = "Công việc hôm qua"
	nowTitle      = "Công việc hôm nay"
	problemTitle  = "Khó khăn"
	solutionTitle = "Giải pháp"
)

func composeSummary(messages []slack.Message, users []slack.User) string {
	var table confluence.Table
	table.Headers = []string{domainTitle, beforeTitle, nowTitle, problemTitle, solutionTitle}
	table.Content = make([][]string, 0, len(messages))

	for _, msg := range messages {
		// preprocess
		input := msg.Text
		input = strings.TrimSpace(input)
		input = titleSentence(input)
		input = slack2confluence(input, users)

		// extract input to report
		report, ok := composeReport(input)
		if !ok {
			continue
		}

		// get display name
		var name string
		for _, user := range users {
			if user.ID == msg.User {
				name = user.Profile.DisplayName
			}
		}
		if len(name) == 0 {
			continue
		}

		table.Content = append(table.Content, []string{
			confluence.FormatBold(name), report.before, report.now, report.problem, report.solution,
		})
	}

	return confluence.ComposeTableFormat(table)
}

// abc -> Abc
// Abc. xyz -> Abc. XYZ
func titleSentence(input string) string {
	// . a -> . A
	regex := regexp.MustCompile(`^.|\.\s.`)
	if !regex.MatchString(input) {
		return input
	}

	replaceFn := func(s string) string {
		return strings.ToUpper(s)
	}

	input = regex.ReplaceAllStringFunc(input, replaceFn)

	return input
}

func slack2confluence(input string, users []slack.User) string {
	input = slack.RemoveBold(input)
	input = slack.NormalizeUser(input, users)
	input = convertSlack2ConfluenceLinks(input)
	input = confluence.ComposeListFormat(input)

	return input
}

func convertSlack2ConfluenceLinks(input string) string {
	links, ok := slack.ExtractLinks(input)
	if !ok {
		return input
	}

	for _, link := range links {
		input = strings.ReplaceAll(input,
			link.Original, confluence.ComposeLinkFormat(link.URL, link.Description))
	}

	return input
}
