package httpwrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddParams(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		parans      []Param
		wantResult  string
		wantError   error
	}{
		{
			name:        "success",
			originalURL: "https://www.google.com",
			parans: []Param{
				{
					Name:  "a",
					Value: "1",
				},
				{
					Name:  "b",
					Value: "2",
				},
			},
			wantResult: "https://www.google.com?a=1&b=2",
			wantError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotResult, gotError := AddParams(tc.originalURL, tc.parans...)
			if tc.wantError != nil {
				assert.Equal(t, tc.wantError, gotError)
				return
			}
			assert.Equal(t, tc.wantResult, gotResult)
		})
	}
}
