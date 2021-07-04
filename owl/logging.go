package owl

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	//nolint:gosimple
	testifyDetection     = regexp.MustCompile("^\\s+Error Trace:")
	testifyMessageMarker = regexp.MustCompile("\n\\s+Messages:")
)

// Errorf is provided for compatibility with testify/require, it will print the errors to stderr.
// Unless verbose is set, testify messages are detected and shortened to the message line only.
func (o *Base) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(strings.TrimPrefix(format, "\n"), args...)
	if !o.Verbose && testifyDetection.MatchString(message) {
		if pos := testifyMessageMarker.FindStringIndex(message); len(pos) == 2 {
			message = message[pos[1]:]
		}
	}
	o.logger.Println(strings.TrimSpace(message))
}

var errFailNow = errors.New("FailNow called by subcommand")

// FailNow is provided for compatibility with testify/require, a panic will trigger an exit with code 1
func (o *Base) FailNow() {
	panic(errFailNow)
}

// Printf wraps fnt.Printf to a configurable stdout, to enable unit testing
func (o *Base) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(o.stdout, format, a...)
}

// Println wraps fnt.Println to a configurable stdout, to enable unit testing
func (o *Base) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(o.stdout, a...)
}
