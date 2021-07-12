package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
	"github.com/xvello/cmds/owl/must"
)

var (
	labelPrefix  = []byte("^XA")
	base64Prefix = []byte("XlhB") // base64 for `^XA`
)

type ZplViewCmd struct {
	Contents string `arg:"positional" help:"ZPL contents to render, can be base64 encoded"`
}

func (c *ZplViewCmd) Run(o owl.Owl) {
	require.NotEmpty(o, c.Contents, "empty zpl contents")

	// Detect and decode base64 encoding
	contents := []byte(c.Contents)
	if bytes.HasPrefix(contents, base64Prefix) {
		var err error
		contents, err = base64.StdEncoding.DecodeString(c.Contents)
		require.NoError(o, err)
	}

	// Look for "begin label" ZPL command
	require.True(o, bytes.HasPrefix(bytes.TrimSpace(contents), labelPrefix), "invalid ZPL data, should start with %s", labelPrefix)

	// Render to multi-page PDF using the labelary web service
	req, err := http.NewRequest(http.MethodPost, "http://api.labelary.com/v1/printers/8dpmm/labels/4x8/", bytes.NewReader(contents))
	require.NoError(o, err)
	req.Header.Add("Accept", "application/pdf")
	res, err := http.DefaultClient.Do(req)
	require.NoError(o, err)
	require.Equal(o, http.StatusOK, res.StatusCode, "unexpected status %d: %s", res.StatusCode, res.Status)

	// Write to file in tmp and open it
	file, err := ioutil.TempFile("", "label-*.pdf")
	require.NoError(o, err)
	_, err = io.Copy(file, res.Body)
	require.NoError(o, err)
	require.NoError(o, file.Close())
	openAndDeleteFile(o, file.Name())
}

// TODO: refactor in owl/must or use pkg/browser
func openAndDeleteFile(o owl.Owl, path string) {
	switch runtime.GOOS {
	case "linux":
		must.Exec(o, "xdg-open", path)
	case "darwin":
		must.Exec(o, "open", path)
	default:
		o.Printf("PDF written to %s", path)
		return // don't delete the file
	}
	time.Sleep(250 * time.Millisecond)
	require.NoError(o, os.Remove(path))
}
