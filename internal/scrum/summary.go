package scrum

import (
	"collector/pkg/confluence"
	"collector/pkg/slack"
	"regexp"
	"strings"
)

const (
	humanMessageIntro   = "Em tổng hợp công việc hôm nay"
	summaryMessageIntro = "Anh nhớ update vào Confluence nhé"

	domainTitle   = "Domain"
	beforeTitle   = "Công việc hôm qua"
	nowTitle      = "Công việc hôm nay"
	problemTitle  = "Khó khăn"
	solutionTitle = "Giải pháp"
	comradeTitle  = "Đồng chí"

	// before, now, problem, solution
	reportNumbers = 4
)

func composeSummary(messages []slack.Message, users []slack.User) (humanSummary []interface{}, confluenceSummary string) {
	humanSummary = make([]interface{}, 0, len(messages)*reportNumbers)
	humanSummary = append(humanSummary, slack.BuildSectionBlock(humanMessageIntro))
	humanSummary = append(humanSummary, slack.BuildDividerBlock())

	var table confluence.Table
	table.Headers = []string{domainTitle, beforeTitle, nowTitle, problemTitle, solutionTitle}
	table.Content = make([][]string, 0, len(messages))

	for _, msg := range messages {
		profile, ok := getProfileOfMessage(msg, users)
		if !ok || profile.DisplayName == "" {
			continue
		}

		// human
		processedMsg := processMessage(msg)
		humanReport, ok := composeReport(processedMsg)
		if !ok {
			continue
		}

		humanMsg := slack.AddBold(comradeTitle) + " " + slack.AddBold(profile.DisplayName) + "\n" +
			slack.AddBold(beforeTitle) + ":" + "\n" +
			humanReport.before + "\n" +
			slack.AddBold(nowTitle) + ":" + "\n" +
			humanReport.now + "\n"

		if humanReport.problem != "" {
			humanMsg += slack.AddBold(problemTitle) + ":" + "\n" +
				humanReport.problem
		}

		if humanReport.solution != "" {
			humanMsg += slack.AddBold(solutionTitle) + ":" + "\n" +
				humanReport.solution
		}

		humanSummary = append(humanSummary, slack.BuildSectionBlock(humanMsg))
		humanSummary = append(humanSummary, slack.BuildDividerBlock())

		// confluence
		confluenceMsg := slack2confluence(processedMsg, users)
		confluenceReport, ok := composeReport(confluenceMsg)
		if !ok {
			continue
		}

		table.Content = append(table.Content, []string{
			confluence.FormatBold(profile.DisplayName),
			confluenceReport.before, confluenceReport.now,
			confluenceReport.problem, confluenceReport.solution,
		})
	}

	confluenceSummary = summaryMessageIntro + "\n" +
		confluence.ComposeTableFormat(table)
	return
}

func processMessage(message slack.Message) string {
	result := message.Text
	result = strings.TrimSpace(result)
	result = titleSentence(result)
	return result
}

func getProfileOfMessage(message slack.Message, users []slack.User) (slack.Profile, bool) {
	for _, user := range users {
		if user.ID == message.User {
			return user.Profile, true
		}
	}

	return slack.Profile{}, false
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
