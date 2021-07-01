package owl

import (
	"log"
	"os"
	"reflect"

	"github.com/alexflint/go-arg"
	"github.com/stretchr/testify/require"
)

// Owl provides helpers to quickly bootstrap a multi-command binary
type Owl struct {
	Aliases *buildAliasesCmd `arg:"subcommand:build-aliases"`
	Verbose bool             `arg:"-v" help:"display full errors"`
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

	// Parse arguments and run a command
	parser := arg.MustParse(cmds)
	if c, ok := parser.Subcommand().(advancedRunnable); ok {
		require.NoError(owl, c.Run(owl, cmds))
	} else if c, ok := parser.Subcommand().(simpleRunnable); ok {
		require.NoError(owl, c.Run(owl))
	} else {
		parser.WriteUsage(os.Stdout)
	}
}
