package confluence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "list",
			input: "- a\n- b\n",
			want:  "* a\n* b\n",
		},
		{
			name:  "list",
			input: "-  a\n - b\n",
			want:  "* a\n* b\n",
		},
		{
			name:  "list",
			input: "-  a\n -  b\n",
			want:  "* a\n* b\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeList(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTitleList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "list",
			input: "* a\n * b\n",
			want:  "* A\n * B\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := titleList(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}
