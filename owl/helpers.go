package owl

import (
	"os/exec"
	"strings"

	"github.com/stretchr/testify/require"
)

// Exec wraps execution of an external command.
// It the command fails, its output is printed and the command stops.
// It the command succeeds, its output is returned as a string.
func (o *Owl) Exec(name string, args ...string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	require.NoError(o, err, string(out))
	return strings.TrimSpace(string(out))
}
