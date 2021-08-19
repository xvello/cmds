package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
	"github.com/xvello/cmds/owl/must"
	"golang.org/x/mod/semver"
)

const kubeNamespace = "bigblue-prod"

var versionPattern = regexp.MustCompile("^v[0-9]+.[0-9]+.[0-9]+$")

type DeployCmd struct {
	Service string `arg:"positional" help:"name of the service to deploy"`
	Version string `arg:"positional" help:"version to deploy"`
	Force   bool   `help:"allow downgrading version"`
}

func (c *DeployCmd) Run(o owl.Owl) {
	require.NotEmpty(o, c.Service, "empty service name")
	require.NotEmpty(o, c.Version, "empty version")
	require.True(o, versionPattern.MatchString(c.Version), "bad version")

	deployName := "bigblue-" + c.Service
	currentImage := must.Exec(o, "kubectl", "--namespace", kubeNamespace, "get", "deployment", deployName, "-o=jsonpath={$.spec.template.spec.containers[0].image}")
	o.Println("Current image: ", currentImage)
	currentParts := strings.Split(currentImage, ":")
	require.Len(o, currentParts, 2, "unexpected image: %s", currentImage)
	require.True(o, semver.IsValid(currentParts[1]), "invalid tag: %s", currentParts[1])

	switch semver.Compare(currentParts[1], c.Version) {
	case -1:
		// All good, new version is higher
	case 0:
		o.Println("this version is already deployed")
		return
	case 1:
		// Downgrading
		if c.Force {
			o.Println("⚠️downgrading to an older version")
		} else {
			require.FailNow(o, "downgrade blocked, use '--force' to allow")
		}
	}

	newImage := fmt.Sprintf("%s:%s", currentParts[0], c.Version)
	o.Println("Deploying: ", newImage)
	must.Exec(o, "kubectl", "--namespace", kubeNamespace, "set", "image", "deployment", deployName, "bigblue="+newImage)

	cmd := exec.Command("kubectl", "--namespace", kubeNamespace, "rollout", "status", "deploy", deployName)
	cmd.Stdout = os.Stdout
	assert.NoError(o, cmd.Run())
}
