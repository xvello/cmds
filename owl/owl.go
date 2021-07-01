package owl

import (
	"io"
	"log"
	"os"
	"reflect"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/require"
)

// Owl provides helpers to quickly bootstrap a multi-command binary
type Owl struct {
	Verbose bool `arg:"-v" help:"display full errors"`

	// Can be overridden for unit tests, defaults to os.StdOut/Err
	stdout io.Writer
	stderr io.Writer
	logger *log.Logger

	// To test FailNow
	mockFailNow      bool
	triggeredFailNow bool
}

type simpleRunnable interface {
	Run(*Owl) error
}

type advancedRunnable interface {
	Run(*Owl, interface{}) error
}

// RunOwl is the entrypoint to call with your command struct.
// The arguments will be parsed and the relevant command called.
func RunOwl(cmds interface{}) {
	// Ensure the given struct embeds Owl and extract it
	owlField := reflect.ValueOf(cmds).Elem().FieldByName("Owl")
	if !owlField.IsValid() {
		log.Fatalf("type %s does not embed the Owl type", reflect.TypeOf(cmds))
	}
	owl, ok := owlField.Addr().Interface().(*Owl)
	if !ok {
		log.Fatalf("Owl field in type %s is not of type %s", reflect.TypeOf(cmds), reflect.TypeOf(new(Owl)))
	}

	// Allow to direct these to buffers for unit tests
	if owl.stderr == nil {
		owl.stderr = os.Stderr
	}
	if owl.stdout == nil {
		owl.stdout = os.Stdout
	}
	if owl.logger == nil {
		owl.logger = log.New(owl.stderr, " ERROR: ", 0)
	}

	// Parse arguments and run a command
	parser := arg.MustParse(cmds)
	if c, ok := parser.Subcommand().(advancedRunnable); ok {
		require.NoError(owl, c.Run(owl, cmds))
	} else if c, ok := parser.Subcommand().(simpleRunnable); ok {
		require.NoError(owl, c.Run(owl))
	} else {
		require.Empty(owl, parser.SubcommandNames(), "command does not implement Run()")
		parser.WriteUsage(os.Stdout)
	}
}
