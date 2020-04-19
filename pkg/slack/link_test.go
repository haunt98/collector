package slack

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestExtractLinks(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantResult []Link
		wantOK     bool
	}{
		{
			name:  "mrkdwn URL and description",
			input: `<http://www.foo.com|This message *is* a link> <www.foo.com|Another link>`,
			wantResult: []Link{
				{
					Original:    "<http://www.foo.com|This message *is* a link>",
					URL:         "http://www.foo.com",
					Description: "This message *is* a link",
				},
				{
					Original:    "<www.foo.com|Another link>",
					URL:         "www.foo.com",
					Description: "Another link",
				},
			},
			wantOK: true,
		},
		{
			name:  "mrkdwn URL",
			input: "<http://www.foo.com> <www.foo.com>",
			wantResult: []Link{
				{
					Original:    "<http://www.foo.com>",
					URL:         "http://www.foo.com",
					Description: "",
				},
				{
					Original:    "<www.foo.com>",
					URL:         "www.foo.com",
					Description: "",
				},
			},
			wantOK: true,
		},
		{
			name:  "URL",
			input: "http://www.foo.com www.foo.com",
			wantResult: []Link{
				{
					Original:    "http://www.foo.com",
					URL:         "http://www.foo.com",
					Description: "",
				},
				{
					Original:    "www.foo.com",
					URL:         "www.foo.com",
					Description: "",
				},
			},
			wantOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotResult, gotOK := ExtractLinks(tc.input)

			assert.Equal(t, tc.wantOK, gotOK)
			assert.Equal(t, tc.wantResult, gotResult)
		})
	}
}
