package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsWorkspaceParamsValid(t *testing.T) {
	tests := []struct {
		testDescription string // Test description
		name            string // Name of workspace
		want            bool
	}{
		{"Normal workspace params", "test", true},
		{"Empty name", "", false},
	}

	for _, test := range tests {
		t.Run(test.testDescription, func(t *testing.T) {
			params := postWorkspaceForm{
				Name: test.name,
			}

			have := isWorkspaceParamsValid(params)

			require.Equal(t, test.want, have)
		})
	}
}
