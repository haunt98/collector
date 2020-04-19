package confluence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeTableFormat(t *testing.T) {
	tests := []struct {
		name  string
		input Table
		want  string
	}{
		{
			name: "only headers",
			input: Table{
				Headers: []string{"a", "b"},
				Content: nil,
			},
			want: "|| a || b ||\n",
		},
		{
			name: "headers and content",
			input: Table{
				Headers: []string{"a", "b"},
				Content: [][]string{
					{"1", "2"},
					{"3", "4"},
				},
			},
			want: "|| a || b ||\n| 1 | 2 |\n| 3 | 4 |\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ComposeTableFormat(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}
