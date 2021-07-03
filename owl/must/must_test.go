package must

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl/mocks"
)

func TestExec(t *testing.T) {
	tests := map[string]struct {
		command     string
		args        []string
		expectedOut string
		expectFail  string
	}{
		"echo_nominal": {
			command:     "echo",
			args:        []string{"one", "two"},
			expectedOut: "one two",
		},
		"echo_split_command": {
			command:     "echo a b c",
			expectedOut: "a b c",
		},
		"unknown": {
			command:    "unknown--command__",
			expectFail: "executable file not found",
		},
		"false": {
			command:    "false",
			expectFail: "exit status 1",
		},
	}
	for name, tc := range tests {
		require.NoError(t, os.Setenv("LANG", "C"))
		t.Run(name, func(t *testing.T) {
			mowl := new(mocks.Owl)
			if tc.expectFail != "" {
				mowl.ExpectRequireFailure(t, tc.expectFail)
				Exec(mowl, tc.command, tc.args...)
			} else {
				assert.Equal(t, tc.expectedOut, Exec(mowl, tc.command, tc.args...))
			}
			mock.AssertExpectationsForObjects(t, mowl)
		})
	}
}
