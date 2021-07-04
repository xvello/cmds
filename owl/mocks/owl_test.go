package mocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
)

const assertionFailureMessage = "assertion failure message"

// TestMockCoverage ensures MockOwl is up to date when the interface is expanded.
func TestMockCoverage(t *testing.T) {
	assert.Implements(t, (*owl.Owl)(nil), new(Owl), "Mock is incomplete, please run `go generate`")
}

func TestExpectAssertFailure(t *testing.T) {
	mowl := new(Owl)
	mowl.ExpectAssertFailure(t, assertionFailureMessage)
	assert.Equal(mowl, 1, 2, assertionFailureMessage)
	mock.AssertExpectationsForObjects(t, mowl)
}

func TestExpectRequireFailure(t *testing.T) {
	mowl := new(Owl)
	mowl.ExpectRequireFailure(t, assertionFailureMessage)
	require.Panics(t, func() {
		require.Equal(mowl, 1, 2, assertionFailureMessage)
	})
	mock.AssertExpectationsForObjects(t, mowl)
}

func TestMock(t *testing.T) {
	mowl := new(Owl)
	mowl.On("IsVerbose").Return(true)
	assert.True(t, mowl.IsVerbose())
	mowl.On("Printf", "format", "arg1", "arg2").Return()
	mowl.Printf("format", "arg1", "arg2")
	mock.AssertExpectationsForObjects(t, mowl)
}
