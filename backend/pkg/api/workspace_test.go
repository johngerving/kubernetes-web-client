package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsWorkspaceParamsValid(t *testing.T) {
	tests := []struct {
		testDescription string            // Test description
		name            string            // Name of workspace
		want            map[string]string // List of problems
	}{
		{"Normal workspace params", "test", map[string]string{}},
		{"Empty name", "", map[string]string{"name": "Name must be at least 2 characters long"}},
	}

	for _, test := range tests {
		t.Run(test.testDescription, func(t *testing.T) {
			params := postWorkspaceForm{
				Name: test.name,
			}

			have := params.valid()

			require.Equal(t, test.want, have)
		})
	}
}
