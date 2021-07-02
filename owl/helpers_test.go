package owl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	tests := map[string]struct {
		command     string
		args        []string
		expectedOut string
		expectFail  bool
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
		"false": {
			command:    "unknown--command__",
			expectFail: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := buildTestCommand()
			c.customFn = func(owl Owl) {
				if len(tc.args) > 0 {
					owl.Printf(owl.Exec(tc.command, tc.args...))
				} else {
					owl.Printf(owl.Exec(tc.command))
				}
			}
			os.Args = []string{"owl", "custom"}
			RunOwl(c)

			assert.Equal(t, tc.expectedOut, c.stdout.String())
			if tc.expectFail {
				assert.True(t, c.triggeredFailNow)
				assert.NotEmpty(t, c.stderr.String())
			} else {
				assert.False(t, c.triggeredFailNow)
				assert.Empty(t, c.stderr.String())
			}
		})
	}
}
