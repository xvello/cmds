package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
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
	githubRemotePattern  = regexp.MustCompile(`git@github.com:([[:word:]]+\/[[:word:]]+).git`)
	githubOpenPRURL      = "https://github.com/%s/compare/%s?expand=1"
)

func (c *NewPrCmd) Run(o owl.Owl) {
	require.NotEmpty(o, c.Branch, "empty branch name")

	// Ensure we have changes to commit
	require.NotEmpty(o, must.Exec(o, "git diff --shortstat HEAD"), "no changes to commit")

	// Ensure we are on the default branch
	currentBranch := must.Exec(o, "git rev-parse --abbrev-ref HEAD")
	if err := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD").Run(); err != nil {
		// Try to recover if the origin head is not already set
		must.Exec(o, "git remote set-head --auto origin")
	}
	defaultBranch := must.Exec(o, "git rev-parse --abbrev-ref origin/HEAD")
	require.Equal(o, strings.Split(defaultBranch, "/")[1], currentBranch, "not on default branch")

	// Validate (and optionally prefix) target branch name
	name := validateBranchName(o, c.Branch)

	o.Printf("Creating and pushing new branch: %s\n", name)
	must.Exec(o, "git", "checkout", "-b", name)
	must.Exec(o, "git", "commit", "-a", "-m", name)

	if c.DryRun {
		o.Println("Dry-run mode: not pushing branch")
		return
	} else {
		o.Println(must.Exec(o, "git", "push", "--set-upstream", "origin", name))
	}

	repository := githubRemotePattern.FindStringSubmatch(must.Exec(o, "git remote get-url origin"))
	if len(repository) == 2 {
		url := fmt.Sprintf(githubOpenPRURL, repository[1], name)
		switch runtime.GOOS {
		case "linux":
			must.Exec(o, "xdg-open", url)
		case "darwin":
			must.Exec(o, "open", url)
		default:
			o.Printf("Open a PR at: %s\n", url)
		}
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
