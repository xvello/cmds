package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ExpectAssertFailure configures the mock to expect a testify/assert failure containing the given string
func (_m *Owl) ExpectAssertFailure(t *testing.T, contains string) {
	_m.On("Errorf", mock.MatchedBy(func(format string) bool {
		return format == "\n%s"
	}), mock.MatchedBy(func(arg string) bool {
		return assert.Contains(t, arg, contains)
	})).Return()
}

// ExpectRequireFailure configures the mock to expect a testify/require failure containing the given string
func (_m *Owl) ExpectRequireFailure(t *testing.T, contains string) {
	_m.ExpectAssertFailure(t, contains)
	_m.On("FailNow").Return()
}
