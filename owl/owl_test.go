package owl

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestCommand() *testCommand {
	var stdout strings.Builder
	var stderr strings.Builder
	return &testCommand{
		Base: Base{
			stdout:           &stdout,
			stderr:           &stderr,
			propagateFailNow: true,
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

func (t *simpleSub) Run(o Owl) {
	o.Println("hello")
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
	require.NotPanics(t, func() { RunOwl(c) })

	require.NotNil(t, c.Simple)
	require.True(t, c.Simple.called)
	require.True(t, c.Simple.Option)
	require.Equal(t, "hello\n", c.stdout.String())
	require.Nil(t, c.Fallible)
	require.False(t, c.passed)
}

func TestFallibleCommand_Ok(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "another", "gopher"}
	require.NotPanics(t, func() { RunOwl(c) })

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
	require.Panics(t, func() { RunOwl(c) })
	require.True(t, strings.HasSuffix(c.stderr.String(), "\tI failed\n"))
}

func TestBadCommand(t *testing.T) {
	c := buildTestCommand()
	os.Args = []string{"owl", "bad"}
	require.Panics(t, func() { RunOwl(c) })
	require.Empty(t, c.stdout.String())
	require.Equal(t, " ERROR: command does not implement Run()\n", c.stderr.String())
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
	require.False(t, c.propagateFailNow)
}

// TestInterfaceCoverage ensures all public methods of Base are exported in the Owl interface.
func TestInterfaceCoverage(t *testing.T) {
	concrete := reflect.TypeOf(new(Base))
	iface := reflect.TypeOf((*Owl)(nil)).Elem()

	// Check that method names match
	for i := 0; i < concrete.NumMethod(); i++ {
		name := concrete.Method(i).Name
		_, found := iface.MethodByName(name)
		assert.True(t, found, "Method %s not in interface", name)
	}

	// Check that method signatures match
	assert.Implements(t, (*Owl)(nil), new(Base), "Method signatures don't match")
}
