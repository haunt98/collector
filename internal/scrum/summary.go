package scrum

import (
	"collector/pkg/slack"
	"fmt"
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
		result += "| " + s.name + " | " + s.before + " | " + s.now + " | " + s.problem + " | " + s.solution + " |\n"
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
		ok, r = makeReport(msg.Text)
		if !ok {
			continue
		}
		r.name = cleanedUsers[msg.User]
		reports = append(reports, r)
	}
	return reports
}

func makeReport(text string) (ok bool, r report) {
	defer func() {
		r.before = beautifyText(r.before)
		r.now = beautifyText(r.now)
		r.problem = beautifyText(r.problem)
		r.solution = beautifyText(r.solution)
	}()

	text = simplyfyText(text)

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

func simplyfyText(text string) string {
	// remove *
	text = strings.ReplaceAll(text, "*", "")

	return text
}

func beautifyText(text string) string {
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
