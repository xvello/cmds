package owl

import (
	"io/fs"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExamples(t *testing.T) {
	// Create a temporary folder to hold binaries
	tmpDir, err := os.MkdirTemp("", "example")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir) // clean up

	// List all examples and pre-compile them
	pwd, err := os.Getwd()
	require.NoError(t, err)
	entries, err := fs.ReadDir(os.DirFS(pwd), "examples")
	require.NoError(t, err)
	for _, e := range entries {
		if e.IsDir() {
			source := path.Join(pwd, "examples", e.Name())
			out, err := exec.Command("go", "build", "-o", tmpDir, source).CombinedOutput()
			require.NoError(t, err, string(out))
		}
	}

	tests := map[string]struct {
		example     string
		args        []string
		expectedOut string
		expectedErr string
	}{
		"calc_usage": {
			example:     "calc",
			expectedOut: "Usage: calc [--verbose] <command>",
		},
		"calc_add_noargs": {
			example:     "calc",
			args:        []string{"add"},
			expectedErr: "numbers is required",
		},
		"calc_add": {
			example:     "calc",
			args:        []string{"add", "1", "2", "3"},
			expectedOut: "6",
		},
		"calc_divide_zero": {
			example:     "calc",
			args:        []string{"divide", "1", "0"},
			expectedErr: "cannot divide by zero",
		},
		"calc_divide_zero_verbose": {
			example:     "calc",
			args:        []string{"divide", "1", "0", "--verbose"},
			expectedErr: "Error Trace:\tmain.go",
		},
		"calc_divide_int": {
			example:     "calc",
			args:        []string{"divide", "100", "10"},
			expectedOut: "10",
		},
		"calc_divide_float": {
			example:     "calc",
			args:        []string{"divide", "10", "4"},
			expectedOut: "2.5",
		},
		"strings_lower": {
			example:     "strings",
			args:        []string{"lower", "COCO"},
			expectedOut: "coco",
		},
		"strings_lower_reverse": {
			example:     "strings",
			args:        []string{"lower", "--reverse", "COCO"},
			expectedOut: "ococ",
		},
		"strings_upper": {
			example:     "strings",
			args:        []string{"upper", "coco"},
			expectedOut: "COCO",
		},
		"strings_upper_reverse": {
			example:     "strings",
			args:        []string{"upper", "--reverse", "CoCo"},
			expectedOut: "OCOC",
		},
	}
	for name, tc := range tests {
		require.NoError(t, os.Setenv("LANG", "C"))
		t.Run(name, func(t *testing.T) {
			out, err := exec.Command(path.Join(tmpDir, tc.example), tc.args...).CombinedOutput()
			if tc.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, string(out), tc.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, string(out), tc.expectedOut)
			}
		})
	}
}
