package owl

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/stretchr/testify/require"
)

type buildAliasesCmd struct{}

func (c *buildAliasesCmd) Run(o *Owl, cmds interface{}) error {
	binary, err := os.Executable()
	require.NoError(o, err, "cannot find current binary path")

	commands := reflect.TypeOf(cmds).Elem()
	for i := 0; i < commands.NumField(); i++ {
		argTags := commands.Field(i).Tag.Get("arg")
		for _, tag := range strings.Split(argTags, ",") {
			if strings.HasPrefix(tag, "subcommand:") {
				name := strings.TrimPrefix(tag, "subcommand:")
				fmt.Printf("alias %s=\"%s %s\"\n", name, binary, name)
			}
		}
	}
	return nil
}
