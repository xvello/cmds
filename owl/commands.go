package owl

import (
	"os"
	"reflect"
	"strings"

	"github.com/stretchr/testify/require"
)

// Extras registers additional subcommands that can be helpful
type Extras struct {
	Aliases *bashAliasesCmd `arg:"subcommand:build-bash-aliases" help:"generate bash aliases for all subcommands"`
}

type bashAliasesCmd struct{}

func (c *bashAliasesCmd) Run(o Owl) error {
	binary, err := os.Executable()
	require.NoError(o, err, "cannot find current binary path")

	commands := reflect.TypeOf(o).Elem()
	for i := 0; i < commands.NumField(); i++ {
		argTags := commands.Field(i).Tag.Get("arg")
		for _, tag := range strings.Split(argTags, ",") {
			if strings.HasPrefix(tag, "subcommand:") {
				name := strings.TrimPrefix(tag, "subcommand:")
				o.Printf("alias %s='%s %s'\n", name, binary, name)
			}
		}
	}
	return nil
}
