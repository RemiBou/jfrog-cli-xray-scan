//+build itest

package main

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"testing"
)

// Smoke test that requires a proper jfrog cli config where it's executed
func Test_Main(t *testing.T) {
	os.Args = []string{"cmd", "s", "-component", "golang.org/x/net v1.8.2"}

	stdout := captureStdout(func() {
		plugins.PluginMain(getApp())
	})

	assert.Matches(t, stdout, "golang.org/x/net:1.8.2")
}

func captureStdout(something func()) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	something()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}
