package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/xvello/owl/mocks"
)

func TestValidateBranchName(t *testing.T) {
	tests := map[string]string{
		"321-a-new-feature":             "feature/big-321-a-new-feature",
		"feature/big-321-a-new-feature": "feature/big-321-a-new-feature",
		"fix/oops-2-things":             "fix/oops-2-things",
		"devops/ci-work":                "devops/ci-work",
		"devop/invalid-prefix":          "",
		"fix/underscores_not_allowed":   "",
		"no-prefix":                     "",
	}
	for input, output := range tests {
		t.Run(input, func(t *testing.T) {
			mowl := new(mocks.Owl)
			if output != "" {
				assert.Equal(t, output, validateBranchName(mowl, input))
			} else {
				mowl.ExpectRequireFailure(t, "invalid branch name "+input)
				assert.Panics(t, func() { validateBranchName(mowl, input) })
			}
			mock.AssertExpectationsForObjects(t, mowl)
		})
	}
}
