package owl

import (
	"fmt"
	"os"
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

// FailNow is provided for compatibility with testify/require, program will exit with code 1
func (o *Base) FailNow() {
	if o.mockFailNow {
		o.triggeredFailNow = true
	} else {
		os.Exit(1)
	}
}

// Printf wraps fnt.Printf to a configurable stdout, to enable unit testing
func (o *Base) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(o.stdout, format, a...)
}
