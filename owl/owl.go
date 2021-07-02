package owl

import (
	"io"
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/require"
)

// Owl provides helpers for your commands, see Base for documentation.
// Commands can cast it to your root type to access global options.
type Owl interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Printf(format string, a ...interface{})
	IsVerbose() bool
	Exec(command string, args ...string) string
	getOwl() *Base
}

// Base provides the base logic to detect and run commands.
// It must be embedded in your root command struct.
type Base struct {
	Verbose bool `arg:"-v" help:"display full errors"`

	// Can be overridden for unit tests, defaults to os.StdOut/Err
	stdout io.Writer
	stderr io.Writer
	logger *log.Logger

	// To test FailNow
	mockFailNow      bool
	triggeredFailNow bool
}

// IsVerbose returns true if the `--verbose` flag has been given.
func (o *Base) IsVerbose() bool {
	return o.Verbose
}

func (o *Base) getOwl() *Base {
	return o
}

type fallibleRunnable interface {
	Run(Owl) error
}

type simpleRunnable interface {
	Run(Owl)
}

// RunOwl is the entrypoint to call with your root struct.
// The arguments will be parsed and the relevant command called.
func RunOwl(root Owl) {
	setupOwl(root)
	parser := arg.MustParse(root)
	if c, ok := parser.Subcommand().(fallibleRunnable); ok {
		require.NoError(root, c.Run(root))
	} else if c, ok := parser.Subcommand().(simpleRunnable); ok {
		c.Run(root)
	} else {
		require.Empty(root, parser.SubcommandNames(), "command does not implement Run()")
		parser.WriteUsage(os.Stdout)
	}
}

// setupOwl sets the required pointer members if they are not set.
// This allows unit tests to override these for behaviour assertions.
func setupOwl(root Owl) {
	owl := root.getOwl()
	if owl.stderr == nil {
		owl.stderr = os.Stderr
	}
	if owl.stdout == nil {
		owl.stdout = os.Stdout
	}
	if owl.logger == nil {
		owl.logger = log.New(owl.stderr, " ERROR: ", 0)
	}
}
