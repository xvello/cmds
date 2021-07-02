package owl

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildTestCommand() *testCommand {
	var stdout strings.Builder
	var stderr strings.Builder
	return &testCommand{
		Base: Base{
			stdout:      &stdout,
			stderr:      &stderr,
			mockFailNow: true,
		},
		stdout: &stdout,
		stderr: &stderr,
	}
}

type testCommand struct {
	Base
	Extras
	Simple   *simpleSub        `arg:"subcommand:simple"`
	Fallible *fallibleSub      `arg:"subcommand:another"`
	Bad      *badSub           `arg:"subcommand:bad"`
	Custom   *customFuncSubCmd `arg:"subcommand:custom"`

	passed   bool
	stdout   *strings.Builder
	stderr   *strings.Builder
	customFn func(Owl)
}

type simpleSub struct {
	Option bool
	called bool
}

func (t *simpleSub) Run(_ Owl) {
	t.called = true
}

type fallibleSub struct {
	Name   string `arg:"positional"`
	Fail   bool
	called bool
}

func (t *fallibleSub) Run(o Owl) error {
	if c, ok := o.(*testCommand); ok {
		c.passed = true
	}
	t.called = true
	if t.Fail {
		return errors.New("I failed")
	}
	return nil
}

type badSub struct{}

type customFuncSubCmd struct{}

func (c *customFuncSubCmd) Run(o Owl) error {
	r, ok := o.(*testCommand)
	require.True(o, ok, "wrong root cmd type")
	r.customFn(o)
	return nil
}

func TestSimpleCommand(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "simple", "--option"}
	RunOwl(c)
	require.False(t, c.triggeredFailNow)
	require.NotNil(t, c.Simple)
	require.True(t, c.Simple.called)
	require.True(t, c.Simple.Option)
	require.Nil(t, c.Fallible)
	require.False(t, c.passed)
}

func TestFallibleCommand_Ok(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "another", "gopher"}
	RunOwl(c)
	require.False(t, c.triggeredFailNow)
	require.NotNil(t, c.Fallible)
	require.True(t, c.Fallible.called)
	require.Equal(t, "gopher", c.Fallible.Name)
	require.True(t, c.passed)
	require.Nil(t, c.Simple)
	require.Empty(t, c.stderr.String())
}

func TestFallibleCommand_Err(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "another", "--fail", "gopher"}
	RunOwl(c)
	require.True(t, c.triggeredFailNow)
	require.True(t, strings.HasSuffix(c.stderr.String(), "\tI failed\n"))
}

func TestBadCommand(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "bad"}
	RunOwl(c)
	require.Empty(t, c.stdout.String())
	require.Equal(t, " ERROR: command does not implement Run()\n", c.stderr.String())
	require.True(t, c.triggeredFailNow)
}

func TestSetupOwl(t *testing.T) {
	c := &struct {
		Base
		Simple *simpleSub `arg:"subcommand:simple"`
	}{}
	os.Args = []string{"owl", "simple"}
	RunOwl(c)
	require.Equal(t, c.stderr, os.Stderr)
	require.Equal(t, c.stdout, os.Stdout)
	require.False(t, c.IsVerbose())
	require.False(t, c.mockFailNow)
	require.False(t, c.triggeredFailNow)
}
