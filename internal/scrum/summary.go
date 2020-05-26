package scrum

import (
	"collector/pkg/confluence"
	"collector/pkg/slack"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	humanMessageIntro = "Em tổng hợp công việc hôm nay " + slack.MentionChannel
	humanEmptyMessage = "Oh, hôm nay không có report"

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

	// only tag last one
	var lastUserID string

	for _, msg := range messages {
		user, ok := getMessageOwner(msg, users)
		if !ok || user.Profile.DisplayName == "" {
			continue
		}

		// human
		processedMsg := processMessage(msg)
		humanReport, ok := composeReport(processedMsg)
		if !ok {
			continue
		}

		humanMsg := slack.AddBold(comradeTitle) + " " + slack.AddBold(user.Profile.DisplayName) + "\n" +
			slack.AddBold(beforeTitle) + ":" + "\n" +
			humanReport.before + "\n" +
			slack.AddBold(nowTitle) + ":" + "\n" +
			humanReport.now + "\n"

		if humanReport.problem != "" {
			humanMsg += slack.AddBold(problemTitle) + ":" + "\n" +
				humanReport.problem + "\n"
		}

		if humanReport.solution != "" {
			humanMsg += slack.AddBold(solutionTitle) + ":" + "\n" +
				humanReport.solution + "\n"
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
			confluence.FormatBold(user.Profile.DisplayName),
			confluenceReport.before, confluenceReport.now,
			confluenceReport.problem, confluenceReport.solution,
		})

		lastUserID = user.ID
	}

	// no comrades -> disable human
	if len(humanSummary) <= 2 {
		humanSummary = append(humanSummary, slack.BuildSectionBlock(humanEmptyMessage))
	}

	// at least 1 comrade report -> enable confluence
	if lastUserID != "" {
		confluenceURL := os.Getenv("CONFLUENCE_URL")
		if confluenceURL == "" {
			log.Fatal("missing confluence link")
		}

		confluenceSummary = fmt.Sprintf("Anh update vào %s nha",
			slack.CreateLink(confluenceURL, "Confluence")) +
			" anh " + slack.MentionUser(lastUserID) + "\n" +
			"```\n" +
			confluence.ComposeTableFormat(table) +
			"```\n"
	}
	return
}

func processMessage(message slack.Message) string {
	result := message.Text
	result = strings.TrimSpace(result)
	result = titleSentence(result)
	return result
}

func getMessageOwner(message slack.Message, users []slack.User) (result slack.User, ok bool) {
	for _, user := range users {
		if user.ID == message.User {
			return user, true
		}
	}

	ok = false
	return
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
