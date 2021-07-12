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

const aliasesPreamble = `#
# Add the following lines to your .bashrc:
# if type -P %[1]s &>/dev/null; then
#     source <(%[1]s build-bash-aliases)
# fi
#
`

type bashAliasesCmd struct{}

func (c *bashAliasesCmd) Run(o Owl) error {
	binary, err := os.Executable()
	require.NoError(o, err, "cannot find current binary path")
	o.Printf(aliasesPreamble, binary)

	commands := reflect.TypeOf(o).Elem()
	for i := 0; i < commands.NumField(); i++ {
		argTags := commands.Field(i).Tag.Get("arg")
		for _, tag := range strings.Split(argTags, ",") {
			if strings.HasPrefix(tag, "subcommand:") {
				name := strings.TrimPrefix(tag, "subcommand:")
				o.Printf("alias %[1]s='%[2]s %[1]s'\n", name, binary)
			}
		}
	}
	return nil
}
