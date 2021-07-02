package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
)

type NewPrCmd struct {
	Branch string `arg:"positional" help:"name of the branch to create"`
	DryRun bool   `arg:"-n" help:"checkout and commit, but don't push'"`
}

func (c *NewPrCmd) Run(o owl.Owl) {
	require.NotEmpty(o, c.Branch, "empty branch name")

	// Ensure we have changes to commit
	require.NotEmpty(o, o.Exec("git diff --shortstat HEAD"), "no changes to commit")

	// Ensure we are on the default branch
	currentBranch := o.Exec("git rev-parse --abbrev-ref HEAD")
	defaultBranch := o.Exec("git rev-parse --abbrev-ref origin/HEAD")
	require.Equal(o, strings.Split(defaultBranch, "/")[1], currentBranch, "not on default branch")

	// Validate and prefix target branch name
	name := c.Branch
	validName := false
	parts := strings.Split(name, "/")
	switch len(parts) {
	case 1:
		if ok, _ := regexp.MatchString("^[0-9]{3,}-[a-z]+", name); ok {
			name = "feature/big-" + name
			validName = true
		}
	case 2:
		switch parts[0] {
		case "feature", "fix", "devops":
			validName = true
		}
	}
	require.True(o, validName, "invalid branch name: %s", name)

	fmt.Printf("Creating and pushing new branch: %s\n", name)
	o.Exec("git", "checkout", "-b", name)
	o.Exec("git", "commit", "-a", "-m", name)

	if c.DryRun {
		fmt.Println("Dry-run mode: not pushing branch")
	} else {
		o.Exec("git push")
	}
}
