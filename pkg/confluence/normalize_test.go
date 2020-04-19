package confluence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeVerticalCharacter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "without link",
			input: "|",
			want:  `\|`,
		},
		{
			name:  "with link",
			input: "| [|]",
			want:  `\| [|]`,
		},
		{
			name:  "with link",
			input: "| [|] [|]",
			want:  `\| [|] [|]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeVerticalCharacter(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNormalizeList(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "list",
			input: "- a\n - b\n",
			want:  "* a\n * b\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeList(tc.input)

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
			got := TitleList(tc.input)

			assert.Equal(t, tc.want, got)
		})
	}
}
