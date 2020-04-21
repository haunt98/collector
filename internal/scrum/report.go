package scrum

import (
	"fmt"
	"regexp"
	"strings"
)

type report struct {
	before, now, problem, solution string
}

type reportType int

const (
	basic reportType = iota + 1
	withProblem
	withSolution
)

func composeReport(input string) (result report, ok bool) {
	withSolutionRegex := composeRegex(withSolution)
	regex := regexp.MustCompile(withSolutionRegex)
	if regex.MatchString(input) {
		submatch := regex.FindStringSubmatch(input)
		for i := range submatch {
			submatch[i] = strings.TrimSpace(submatch[i])
		}

		result.before, result.now, result.problem, result.solution = submatch[1], submatch[2], submatch[3], submatch[4]
		ok = true
		return
	}

	withProblemRegex := composeRegex(withProblem)
	regex = regexp.MustCompile(withProblemRegex)
	if regex.MatchString(input) {
		submatch := regex.FindStringSubmatch(input)
		for i := range submatch {
			submatch[i] = strings.TrimSpace(submatch[i])
		}

		result.before, result.now, result.problem = submatch[1], submatch[2], submatch[3]
		ok = true
		return
	}

	basicRegex := composeRegex(basic)
	regex = regexp.MustCompile(basicRegex)
	if regex.MatchString(input) {
		submatch := regex.FindStringSubmatch(input)
		for i := range submatch {
			submatch[i] = strings.TrimSpace(submatch[i])
		}

		result.before, result.now = submatch[1], submatch[2]
		ok = true
		return
	}

	ok = false
	return
}

var beforeRegexList = []string{
	`yesterday`,
	`h[oô]m\s+qua`,
	`h[oô]m\s+kia`,
	`h[oô]m\s+b[uưữ]a`,
	`h[oô]m\s+tr[uư][oơớ]c`,
	`tu[aâầ]n\s+qua`,
	`tu[aâầ]n\s+tr[uư][ơớ]c`,
	`tu[aâầ]n\s+kia`,
}

var nowRegexList = []string{
	`today`,
	`h[oô]m\s+nay`,
	`tu[aâầ]n\sn[aà]y`,
}

var problemRegexList = []string{
	`problem`,
	`kh[oó]\s+kh[aă]n`,
	`v[aâấ]n\s+[dđ][eêề]`,
}

var solutionRegexList = []string{
	`solution`,
	`gi[aả]i\s+ph[aá]p`,
	`gi[aả]i\s+quy[eêế]t`,
}

func composeRegex(t reportType) string {
	beforePhrase := composeOrRegex(beforeRegexList)
	nowPhrase := composeOrRegex(nowRegexList)
	problemPhrase := composeOrRegex(problemRegexList)
	solutionPhrase := composeOrRegex(solutionRegexList)

	switch t {
	case basic:
		return fmt.Sprintf(`(?is)%s:?(.+?)%s:?(.+)`,
			beforePhrase, nowPhrase)
	case withProblem:
		return fmt.Sprintf(`(?is)%s:?(.+?)%s:?(.+?)%s:?(.+)`,
			beforePhrase, nowPhrase, problemPhrase)
	case withSolution:
		return fmt.Sprintf(`(?is)%s:?(.+?)%s:?(.+?)%s:?(.+?)%s:?(.+)`,
			beforePhrase, nowPhrase, problemPhrase, solutionPhrase)
	default:
		return ""
	}
}

// (?:a|b)
func composeOrRegex(arr []string) string {
	s := ``
	for i := range arr {
		s += arr[i]
		if i != len(arr)-1 {
			s += `|`
		}
	}

	return fmt.Sprintf(`(?:%s)`, s)
}
