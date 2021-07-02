package owl

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorf(t *testing.T) {
	tests := map[string]struct {
		fn             func(Owl)
		verbose        bool
		expected       string
		expectedPrefix string
		expectedSuffix string
		expectFailNow  bool
	}{
		"simple_one_line": {
			fn:       func(o Owl) { o.Errorf("simple error") },
			verbose:  false,
			expected: " ERROR: simple error\n",
		},
		"simple_multi_line": {
			fn:       func(o Owl) { o.Errorf("simple error\non several\nlines") },
			verbose:  false,
			expected: " ERROR: simple error\non several\nlines\n",
		},
		"require_with_message_default": {
			fn:            func(o Owl) { require.Equal(o, 1, 2, "oops") },
			verbose:       false,
			expected:      " ERROR: oops\n",
			expectFailNow: true,
		},
		"assert_with_message_verbose": {
			fn:             func(o Owl) { assert.Equal(o, 1, 2, "oops") },
			verbose:        true,
			expectedPrefix: " ERROR: Error Trace:\tlogging_test.go:",
			expectedSuffix: "Error:      \tNot equal: \n\t            \texpected: 1\n\t            \tactual  : 2\n\tMessages:   \toops\n",
		},
		"assert_no_message": {
			fn:             func(o Owl) { assert.Equal(o, 1, 2) },
			verbose:        false,
			expectedPrefix: " ERROR: Error Trace:\tlogging_test.go:",
			expectedSuffix: "Error:      \tNot equal: \n\t            \texpected: 1\n\t            \tactual  : 2\n",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := buildTestCommand()
			c.customFn = tc.fn

			os.Args = []string{"owl", "custom"}
			if tc.verbose {
				os.Args = append(os.Args, "--verbose")
			}
			RunOwl(c)

			assert.Empty(t, c.stdout.String())
			assert.Equal(t, tc.expectFailNow, c.triggeredFailNow)
			if tc.expected != "" {
				assert.Equal(t, tc.expected, c.stderr.String())
			}
			if tc.expectedPrefix != "" {
				assert.Equal(t, tc.expectedPrefix, c.stderr.String()[0:len(tc.expectedPrefix)])
			}
			if tc.expectedSuffix != "" {
				suffix := c.stderr.String()[c.stderr.Len()-len(tc.expectedSuffix):]
				assert.Equal(t, tc.expectedSuffix, suffix)
			}
		})
	}
}
