package scrum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeReport(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantReport report
		wantOK     bool
	}{
		{
			name:  "basic",
			input: `yesterday: eat today: sleep`,
			wantReport: report{
				before:   "eat",
				now:      "sleep",
				problem:  "",
				solution: "",
			},
			wantOK: true,
		},
		{
			name:  "with problem",
			input: "yesterday: eat today: sleep problem: none",
			wantReport: report{
				before:   "eat",
				now:      "sleep",
				problem:  "none",
				solution: "",
			},
			wantOK: true,
		},
		{
			name:  "with solution",
			input: "yesterday: eat today: sleep problem: none solution keep it going",
			wantReport: report{
				before:   "eat",
				now:      "sleep",
				problem:  "none",
				solution: "keep it going",
			},
			wantOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotReport, gotOK := composeReport(tc.input)

			assert.Equal(t, tc.wantOK, gotOK)
			assert.Equal(t, tc.wantReport, gotReport)
		})
	}
}

func TestComposeOrRegex(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name: "1 element",
			input: []string{
				`a`,
			},
			want: `(?:a)`,
		},
		{
			name: "2 elements",
			input: []string{
				`a`, `b`,
			},
			want: `(?:a|b)`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := composeOrRegex(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}
