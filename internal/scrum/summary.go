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

func makeConfluenceSummary(messages []slack.Message, users []slack.User) string {
	cleanedUsers := cleanUsers(users)
	cleanedMessages := cleanMessages(messages, cleanedUsers)

	result := fmt.Sprintf("|| %s || %s || %s || %s || %s ||\n", domainTitle, beforeTitle, nowTitle, problemTitle, solutionTitle)
	reports := makeReports(cleanedMessages, cleanedUsers)
	for _, report := range reports {
		result += fmt.Sprintf("| *%s* | %s | %s | %s | %s |\n",
			report.name, report.before, report.now, report.problem, report.solution)
	}
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
	text = convertSlack2ConfluenceLists(text)
	text = convertSlackUsers(text, cleanedUsers)
	text = removeVertical(text)

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
	regex := regexp.MustCompile(`(?is)(?:h[oô]m\s+qua|h[oô]m\s+kia|h[oô]m\s+b[uưữ]a|h[oô]m\s+tr[uư][oơớ]c|tu[aâầ]n\s+qua|tu[aâầ]n\s+tr[uư][ơớ]c|tu[aâầ]n\s+kia)(.+?)(?:h[oô]m\s+nay|tu[aâầ]n\sn[aà]y)(.+?)(?:kh[oó]\s+kh[aă]n|v[aâấ]n\s+[dđ][eêề])(.+?)(?:gi[aả]i\s+ph[aá]p|gi[aả]i\s+quy[eêế]t)(.+)`)
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
	regex := regexp.MustCompile(`(?is)(?:h[oô]m\s+qua|h[oô]m\s+kia|h[oô]m\s+b[uưữ]a|h[oô]m\s+tr[uư][oơớ]c|tu[aâầ]n\s+qua|tu[aâầ]n\s+tr[uư][ơớ]c|tu[aâầ]n\s+kia)(.+?)(?:h[oô]m\s+nay|tu[aâầ]n\sn[aà]y)(.+?)(?:kh[oó]\s+kh[aă]n|v[aâấ]n\s+[dđ][eêề])(.+)`)
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
	regex := regexp.MustCompile(`(?is)(?:h[oô]m\s+qua|h[oô]m\s+kia|h[oô]m\s+b[uưữ]a|h[oô]m\s+tr[uư][oơớ]c|tu[aâầ]n\s+qua|tu[aâầ]n\s+tr[uư][ơớ]c|tu[aâầ]n\s+kia)(.+?)(?:h[oô]m\s+nay|tu[aâầ]n\sn[aà]y)(.+)`)
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
	text = strings.ReplaceAll(text, "*", "")
	return text
}

// remove | for not messing with confluence links
func removeVertical(text string) string {
	text = strings.ReplaceAll(text, "|", `\|`)
	return text
}

func convertSlack2ConfluenceLinks(text string) string {
	regex := regexp.MustCompile(`<(http.+)\|.*>`)
	subs := regex.FindAllStringSubmatch(text, -1)
	for _, sub := range subs {
		original := sub[0]
		confluenceLink := "[" + sub[1] + "]"

		text = strings.ReplaceAll(text, original, confluenceLink)
	}

	regex = regexp.MustCompile(`<(http.+)>`)
	subs = regex.FindAllStringSubmatch(text, -1)
	for _, sub := range subs {
		original := sub[0]
		confluenceLink := "[" + sub[1] + "]"

		text = strings.ReplaceAll(text, original, confluenceLink)
	}

	return text
}

func convertSlack2ConfluenceLists(text string) string {
	text = strings.ReplaceAll(text, "• ", "* ")
	text = strings.ReplaceAll(text, "- ", "* ")
	text = strings.ReplaceAll(text, "* * ", "** ")
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
