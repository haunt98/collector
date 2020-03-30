package report

import (
	"collector/pkg/slack"
	"regexp"
	"strings"
)

func MakeMessage(messages []slack.Message, users []slack.User) string {
	cleanedUsers := cleanUsers(users)
	cleanedMessages := cleanMessages(messages, cleanedUsers)

	result := "```\n"
	result += "|| Domain || Công việc hôm qua || Công việc hôm nay || Khó khăn || Giải pháp ||\n"
	multiReport := makeMultiReport(cleanedMessages, cleanedUsers)
	for _, s := range multiReport {
		result += "| " + s.name + " | " + s.yesterday + " | " + s.today + " | " + s.problem + " | " + s.solution + " |\n"
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

type singleReport struct {
	name, yesterday, today, problem, solution string
}

func makeMultiReport(cleanedMessages []slack.Message, cleanedUsers map[string]string) []singleReport {
	multiReport := make([]singleReport, 0, len(cleanedMessages))
	for _, msg := range cleanedMessages {
		var s singleReport
		var ok bool
		ok, s.yesterday, s.today, s.problem, s.solution = readReportFromMessage(msg.Text)
		if !ok {
			continue
		}
		s.name = cleanedUsers[msg.User]
		multiReport = append(multiReport, s)
	}
	return multiReport
}

const (
	yesterdayVie = "hôm qua"
	todayVie     = "hôm nay"
	problemVie   = "khó khăn"
	solutionVie  = "giải pháp"
)

func readReportFromMessage(text string) (ok bool, yesterday, today, problem, solution string) {
	defer func() {
		yesterday = beautifyText(yesterday)
		today = beautifyText(today)
		problem = beautifyText(problem)
		solution = beautifyText(solution)
	}()

	lowerText := strings.ToLower(text)

	yesterdayTodayValid := regexp.MustCompile(`(?s).*` + yesterdayVie + `.*` + todayVie + `.*`)
	if !yesterdayTodayValid.MatchString(lowerText) {
		ok = false
		return
	}

	yesterdayIndex := strings.Index(lowerText, yesterdayVie)
	todayIndex := strings.Index(lowerText, todayVie)

	problemValid := regexp.MustCompile(`(?s).*` + problemVie + `.*`)
	if !problemValid.MatchString(lowerText) {
		ok = true
		yesterday = text[yesterdayIndex+len(yesterdayVie) : todayIndex]
		today = text[todayIndex+len(todayVie):]
		return
	}

	problemIndex := strings.Index(lowerText, problemVie)

	solutionValid := regexp.MustCompile(`(?s).*` + solutionVie + `.*`)
	if !solutionValid.MatchString(lowerText) {
		ok = true
		yesterday = text[yesterdayIndex+len(yesterdayVie) : todayIndex]
		today = text[todayIndex+len(todayVie) : problemIndex]
		problem = text[problemIndex+len(problemVie):]
		return
	}

	solutionIndex := strings.Index(lowerText, solutionVie)

	ok = true
	yesterday = text[yesterdayIndex+len(yesterdayVie) : todayIndex]
	today = text[todayIndex+len(todayVie) : problemIndex]
	problem = text[problemIndex+len(problemVie) : solutionIndex]
	solution = text[solutionIndex+len(solutionVie):]
	return
}

func beautifyText(text string) string {
	if len(text) == 0 {
		return ""
	}

	if text[0] == ':' {
		text = text[1:]
	}

	text = strings.TrimSpace(text)

	return text
}
