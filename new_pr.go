package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
	"github.com/xvello/cmds/owl/must"
)

type NewPrCmd struct {
	Branch string `arg:"positional" help:"name of the branch to create"`
	DryRun bool   `arg:"-n" help:"checkout and commit, but don't push'"`
}

var (
	linearBranchPattern  = regexp.MustCompile("^[0-9]{3,}-[a-z][a-z0-9-]+$")
	genericBranchPattern = regexp.MustCompile("^(feature|fix|devops)/[a-z][a-z0-9-]+$")
)

func (c *NewPrCmd) Run(o owl.Owl) {
	require.NotEmpty(o, c.Branch, "empty branch name")

	// Ensure we have changes to commit
	require.NotEmpty(o, must.Exec(o, "git diff --shortstat HEAD"), "no changes to commit")

	// Ensure we are on the default branch
	currentBranch := must.Exec(o, "git rev-parse --abbrev-ref HEAD")
	defaultBranch := must.Exec(o, "git rev-parse --abbrev-ref origin/HEAD")
	require.Equal(o, strings.Split(defaultBranch, "/")[1], currentBranch, "not on default branch")

	// Validate (and optionally prefix) target branch name
	name := validateBranchName(o, c.Branch)

	fmt.Printf("Creating and pushing new branch: %s\n", name)
	must.Exec(o, "git", "checkout", "-b", name)
	must.Exec(o, "git", "commit", "-a", "-m", name)
	must.Exec(o, "git", "branch", "--set-upstream-to=origin/"+name, name)

	if c.DryRun {
		fmt.Println("Dry-run mode: not pushing branch")
	} else {
		o.Printf(must.Exec(o, "git push"))
	}
}

func validateBranchName(o owl.Owl, name string) string {
	if linearBranchPattern.MatchString(name) {
		name = "feature/big-" + name
	}
	require.True(o, genericBranchPattern.MatchString(name),
		"invalid branch name %s, must match %s", name, genericBranchPattern.String())
	return name
}
