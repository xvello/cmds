package owl

import (
	"os"
	"reflect"
	"strings"

	"github.com/stretchr/testify/require"
)

// ExtraCommands registers additional subcommands that can be helpful:
type ExtraCommands struct {
	Aliases *bashAliasesCmd `arg:"subcommand:bash-aliases" help:"generate bash aliases for all subcommands"`
}

type bashAliasesCmd struct{}

func (c *bashAliasesCmd) Run(o *Owl, cmds interface{}) error {
	binary, err := os.Executable()
	require.NoError(o, err, "cannot find current binary path")

	commands := reflect.TypeOf(cmds).Elem()
	for i := 0; i < commands.NumField(); i++ {
		argTags := commands.Field(i).Tag.Get("arg")
		for _, tag := range strings.Split(argTags, ",") {
			if strings.HasPrefix(tag, "subcommand:") {
				name := strings.TrimPrefix(tag, "subcommand:")
				o.printf("alias %s=\"%s %s\"\n", name, binary, name)
			}
		}
	}
	return nil
}
