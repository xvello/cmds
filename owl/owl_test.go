package owl

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildTestCommand() *testCommand {
	var stdout strings.Builder
	var stderr strings.Builder
	return &testCommand{
		Owl: Owl{
			stdout:      &stdout,
			stderr:      &stderr,
			mockFailNow: true,
		},
		stdout: &stdout,
		stderr: &stderr,
	}
}

type testCommand struct {
	Owl
	ExtraCommands
	Simple   *simpleSub        `arg:"subcommand:simple"`
	Advanced *advancedSub      `arg:"subcommand:another"`
	Bad      *badSub           `arg:"subcommand:bad"`
	Custom   *customFuncSubCmd `arg:"subcommand:custom"`

	passed   bool
	stdout   *strings.Builder
	stderr   *strings.Builder
	customFn func(*Owl)
}

type simpleSub struct {
	Option bool
	called bool
}

func (t *simpleSub) Run(_ *Owl) error {
	t.called = true
	return nil
}

type advancedSub struct {
	Name   string `arg:"positional"`
	called bool
}

func (t *advancedSub) Run(_ *Owl, cmds interface{}) error {
	if c, ok := cmds.(*testCommand); ok {
		c.passed = true
	}
	t.called = true
	return nil
}

type badSub struct{}

type customFuncSubCmd struct{}

func (c *customFuncSubCmd) Run(o *Owl, root interface{}) error {
	r, ok := root.(*testCommand)
	require.True(o, ok, "wrong root cmd type")
	r.customFn(o)
	return nil
}

func TestSimpleCommand(t *testing.T) {
	c := &testCommand{}
	os.Args = []string{"owl", "simple", "--option"}
	RunOwl(c)
	require.NotNil(t, c.Simple)
	require.True(t, c.Simple.called)
	require.True(t, c.Simple.Option)
	require.Nil(t, c.Advanced)
	require.False(t, c.passed)
}

func TestAdvancedCommand(t *testing.T) {
	c := &testCommand{}
	os.Args = []string{"owl", "another", "gopher"}
	RunOwl(c)
	require.NotNil(t, c.Advanced)
	require.True(t, c.Advanced.called)
	require.Equal(t, "gopher", c.Advanced.Name)
	require.True(t, c.passed)
	require.Nil(t, c.Simple)
}

func TestBadCommand(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "bad"}
	RunOwl(c)
	require.Empty(t, c.stdout.String())
	require.Equal(t, " ERROR: command does not implement Run()\n", c.stderr.String())
	require.True(t, c.triggeredFailNow)
}
