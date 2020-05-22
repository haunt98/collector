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

	// before, now, problem, solution
	reportNumbers = 4
)

func composeSummaryForHuman(messages []slack.Message, users []slack.User) []interface{} {
	blocks := make([]interface{}, 0, len(messages)*reportNumbers)

	for _, msg := range messages {
		input := processMessage(msg)

		// extract input to report
		report, ok := composeReport(input)
		if !ok {
			continue
		}

		// get display name
		name := getDisplayName(msg, users)
		if name == "" {
			continue
		}

		// build blocks
		nameBlock := slack.BuildSectionBlock(slack.AddBold(name))
		blocks = append(blocks, nameBlock)

		beforeTitleBlock := slack.BuildSectionBlock(slack.AddBold(beforeTitle))
		blocks = append(blocks, beforeTitleBlock)

		beforeContentBlock := slack.BuildSectionBlock(report.before)
		blocks = append(blocks, beforeContentBlock)

		nowTitleBlock := slack.BuildSectionBlock(slack.AddBold(nowTitle))
		blocks = append(blocks, nowTitleBlock)

		nowContentBlock := slack.BuildSectionBlock(report.now)
		blocks = append(blocks, nowContentBlock)

		if report.problem != "" {
			problemTitleBlock := slack.BuildSectionBlock(slack.AddBold(problemTitle))
			blocks = append(blocks, problemTitleBlock)

			problemContentBlock := slack.BuildSectionBlock(report.problem)
			blocks = append(blocks, problemContentBlock)
		}

		if report.solution != "" {
			solutionTitleBlock := slack.BuildSectionBlock(slack.AddBold(solutionTitle))
			blocks = append(blocks, solutionTitleBlock)

			solutionContentBlock := slack.BuildSectionBlock(report.solution)
			blocks = append(blocks, solutionContentBlock)
		}
	}

	return blocks
}

func composeSummaryForConfluence(messages []slack.Message, users []slack.User) string {
	var table confluence.Table
	table.Headers = []string{domainTitle, beforeTitle, nowTitle, problemTitle, solutionTitle}
	table.Content = make([][]string, 0, len(messages))

	for _, msg := range messages {
		input := processMessage(msg)

		input = slack2confluence(input, users)

		// extract input to report
		report, ok := composeReport(input)
		if !ok {
			continue
		}

		// get display name
		name := getDisplayName(msg, users)
		if name == "" {
			continue
		}

		table.Content = append(table.Content, []string{
			confluence.FormatBold(name), report.before, report.now, report.problem, report.solution,
		})
	}

	return confluence.ComposeTableFormat(table)
}

func processMessage(message slack.Message) string {
	result := message.Text
	result = strings.TrimSpace(result)
	result = titleSentence(result)
	return result
}

func getDisplayName(message slack.Message, users []slack.User) string {
	for _, user := range users {
		if user.ID == message.User {
			return user.Profile.DisplayName
		}
	}

	return ""
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
