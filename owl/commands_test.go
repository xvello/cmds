package owl

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashAliases(t *testing.T) {
	executable, _ := os.Executable()
	expected := fmt.Sprintf("alias simple=\"%s simple\"\nalias another=\"%s another\"\n", executable, executable)
	var stdout strings.Builder
	var stderr strings.Builder
	c := &struct {
		Owl
		ExtraCommands
		Simple   *simpleSub   `arg:"subcommand:simple"`
		Advanced *advancedSub `arg:"subcommand:another"`
	}{Owl: Owl{
		stdout: &stdout,
		stderr: &stderr,
	}}

	os.Args = []string{"owl", "bash-aliases"}
	RunOwl(c)
	assert.Empty(t, stderr.String())
	assert.Equal(t, expected, stdout.String())
}
