package scrum

import (
	"collector/pkg/slack"
	"fmt"
	"log"
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

func makeSummary(messages []slack.Message, users []slack.User) string {
	cleanedUsers := cleanUsers(users)
	cleanedMessages := cleanMessages(messages, cleanedUsers)

	result := "```\n"
	result += fmt.Sprintf("|| %s || %s || %s || %s || %s ||\n", domainTitle, beforeTitle, nowTitle, problemTitle, solutionTitle)
	reports := makeReports(cleanedMessages, cleanedUsers)
	for _, s := range reports {
		result += fmt.Sprintf("| *%s* | %s | %s | %s | %s |\n",
			s.name, s.before, s.now, s.problem, s.solution)
	}
	result += "```"
	return result
}

func cleanUsers(users []slack.User) map[string]string {
	cleanedUsers := make(map[string]string)
	for _, user := range users {
		if user.IsBot || user.IsAppUser {
			continue
		}
		cleanedUsers[user.ID] = user.Profile.DisplayName
	}
	return cleanedUsers
}

func cleanMessages(messages []slack.Message, cleanedUsers map[string]string) []slack.Message {
	cleanedMessages := make([]slack.Message, 0, len(messages))
	for _, msg := range messages {
		if _, ok := cleanedUsers[msg.User]; !ok {
			continue
		}

		cleanedMessages = append(cleanedMessages, msg)
	}
	return cleanedMessages
}

type report struct {
	name, before, now, problem, solution string
}

func makeReports(cleanedMessages []slack.Message, cleanedUsers map[string]string) []report {
	reports := make([]report, 0, len(cleanedMessages))
	for _, msg := range cleanedMessages {
		var r report
		var ok bool
		ok, r = makeReport(msg.Text, cleanedUsers)
		if !ok {
			continue
		}
		r.name = cleanedUsers[msg.User]
		reports = append(reports, r)
	}
	return reports
}

func makeReport(text string, cleanedUsers map[string]string) (ok bool, r report) {
	defer func() {
		r.before = trimSpace(r.before)
		r.now = trimSpace(r.now)
		r.problem = trimSpace(r.problem)
		r.solution = trimSpace(r.solution)
	}()

	text = removeStar(text)
	text = convertSlack2ConfluenceLinks(text)
	text = convertSlackUsers(text, cleanedUsers)

	ok, r = consume4(text)
	if ok {
		return
	}

	ok, r = consume3(text)
	if ok {
		return
	}

	ok, r = consume2(text)
	if ok {
		return
	}

	ok = false
	return
}

func consume4(text string) (ok bool, r report) {
	regex := regexp.MustCompile(`(?is)(?:hôm\s+qua|hôm\s+kia|hôm\s+bữa|hôm\s+trước|tuần\s+qua|tuần\s+trước|tuần\s+kia)(.+?)(?:hôm\s+nay|tuần\snày)(.+?)(?:khó\s+khăn|vấn\s+đề)(.+?)(?:giải\s+pháp|giải\s+quyết)(.+)`)
	if !regex.MatchString(text) {
		ok = false
		return
	}

	subs := regex.FindStringSubmatch(text)
	r.before, r.now, r.problem, r.solution = subs[1], subs[2], subs[3], subs[4]
	ok = true
	return
}

func consume3(text string) (ok bool, r report) {
	regex := regexp.MustCompile(`(?is)(?:hôm\s+qua|hôm\s+kia|hôm\s+bữa|hôm\s+trước|tuần\s+qua|tuần\s+trước|tuần\s+kia)(.+?)(?:hôm\s+nay|tuần\snày)(.+?)(?:khó\s+khăn|vấn\s+đề)(.+)`)
	if !regex.MatchString(text) {
		ok = false
		return
	}

	subs := regex.FindStringSubmatch(text)
	r.before, r.now, r.problem = subs[1], subs[2], subs[3]
	ok = true
	return
}

func consume2(text string) (ok bool, r report) {
	regex := regexp.MustCompile(`(?is)(?:hôm\s+qua|hôm\s+kia|hôm\s+bữa|hôm\s+trước|tuần\s+qua|tuần\s+trước|tuần\s+kia)(.+?)(?:hôm\s+nay|tuần\snày)(.+)`)
	if !regex.MatchString(text) {
		ok = false
		return
	}

	subs := regex.FindStringSubmatch(text)
	r.before, r.now = subs[1], subs[2]
	ok = true
	return
}

// remove *
func removeStar(text string) string {
	return strings.ReplaceAll(text, "*", "")
}

func convertSlack2ConfluenceLinks(text string) string {
	regex := regexp.MustCompile(`\(.*<(http.+)\|.*>.*\)`)
	if !regex.MatchString(text) {
		return text
	}

	subs := regex.FindAllStringSubmatch(text, -1)
	for _, sub := range subs {
		original := sub[0]
		confluenceLink := "[" + sub[1] + "]"

		log.Println("original", original)
		log.Println("confluenceLink", confluenceLink)

		text = strings.ReplaceAll(text, original, confluenceLink)
	}

	return text
}

func convertSlack2ConfluenceLists(text string) string {
	text = strings.ReplaceAll(text, "•", "*")
	text = strings.ReplaceAll(text, "-", "*")
	text = strings.ReplaceAll(text, "* *", "**")
	return text
}

func convertSlackUsers(text string, cleanedUsers map[string]string) string {
	for id, name := range cleanedUsers {
		text = strings.ReplaceAll(text, fmt.Sprintf("<@%s>", id), fmt.Sprintf("@%s", name))
	}
	return text
}

func trimSpace(text string) string {
	if len(text) == 0 {
		return ""
	}

	text = strings.TrimSpace(text)

	if text[0] == ':' {
		text = text[1:]
	}

	text = strings.TrimSpace(text)

	return text
}
