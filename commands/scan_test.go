package commands

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_Mvn_Component_Flag_Scanner(t *testing.T) {
	context := &fakeContext{flags: map[string]string{componentFlagKey: "org.slf4j:slf4j-ext:jar:1.7.26:compile"}}
	fakeXrayClient := &fakeXrayClient{
		resultOK: &ComponentSummaryResult{},
	}
	err := scanCmd(context, os.Stdin, fakeXrayClient.scan)
	require.NoError(t, err)
	require.ElementsMatch(t, []component{"gav:org.slf4j:slf4j-ext:1.7.26"}, fakeXrayClient.scanned)
}

func Test_Mvn_Component_Stdin_Scanner(t *testing.T) {
	context := &fakeContext{flags: map[string]string{}}
	stdin := bytes.NewBufferString(
		"[INFO]    org.slf4j:slf4j-ext:jar:1.7.26:compile -- module slf4j.ext (auto)")
	fakeXrayClient := &fakeXrayClient{
		resultOK: &ComponentSummaryResult{},
	}
	err := scanCmd(context, stdin, fakeXrayClient.scan)
	require.NoError(t, err)
	require.ElementsMatch(t, []component{"gav:org.slf4j:slf4j-ext:1.7.26"}, fakeXrayClient.scanned)
}

// TODO: Get input from files
func Test_Mvn_Dependencies_Stdin_Scanner(t *testing.T) {
	lines := make(chan string)
	go func() {
		lines <- "[INFO]    org.slf4j:slf4j-ext:jar:1.7.26:compile -- module slf4j.ext (auto)"
		lines <- "[INFO]    org.jolokia:jolokia-core:jar:1.6.2:compile -- module jolokia.core (auto)"
		lines <- "[INFO]    com.googlecode.json-simple:json-simple:jar:1.1.1:compile -- module json.simple (auto)"
		lines <- "[INFO]    io.jaegertracing:jaeger-client:jar:1.2.0:compile -- module jaeger.client (auto)"
		lines <- "[INFO]    io.jaegertracing:jaeger-thrift:jar:1.2.0:compile -- module jaeger.thrift (auto)"
		lines <- "[INFO]    org.apache.thrift:libthrift:jar:0.13.0:compile -- module libthrift (auto)"
		lines <- "[INFO]    com.google.code.findbugs:jsr305:jar:2.0.0:compile -- module jsr305 (auto)"
		lines <- "[INFO]    "
		lines <- "[INFO]    "
		lines <- "[INFO] ----------------< org.jfrog.access:access-coverage-all >----------------"
		lines <- "[INFO] Building JFrog Access OSS Coverage 7.x-SNAPSHOT                  [21/21]"
		lines <- "[INFO] --------------------------------[ jar ]---------------------------------"
		lines <- "[INFO]    "
		close(lines)
	}()
	fakeXrayClient := &fakeXrayClient{
		resultOK: &ComponentSummaryResult{},
	}
	err := scan(lines, fakeXrayClient.scan)
	require.NoError(t, err)
	require.ElementsMatch(t, []component{
		"gav:org.jolokia:jolokia-core:1.6.2",
		"gav:org.slf4j:slf4j-ext:1.7.26",
		"gav:com.googlecode.json-simple:json-simple:1.1.1",
		"gav:io.jaegertracing:jaeger-client:1.2.0",
		"gav:io.jaegertracing:jaeger-thrift:1.2.0",
		"gav:com.google.code.findbugs:jsr305:2.0.0",
		"gav:org.apache.thrift:libthrift:0.13.0"}, fakeXrayClient.scanned)

}

func Test_Go_List_Stdin_Scanner(t *testing.T) {
	lines := make(chan string)
	go func() {
		lines <- "github.com/jfrog/jfrog-cli-plugin-template"
		lines <- "github.com/BurntSushi/toml v0.3.1"
		close(lines)
	}()
	fakeXrayClient := &fakeXrayClient{
		resultOK: &ComponentSummaryResult{},
	}
	err := scan(lines, fakeXrayClient.scan)
	require.NoError(t, err)
	require.ElementsMatch(t, []component{
		"go://github.com/BurntSushi/toml:0.3.1",
	}, fakeXrayClient.scanned)
}

func TestValidateOutput(t *testing.T) {

}

func TestOnScanError(t *testing.T) {

}

type fakeXrayClient struct {
	scanned   []component
	resultOK  *ComponentSummaryResult
	resultErr error
}

func (x *fakeXrayClient) scan(comps []component) (*ComponentSummaryResult, error) {
	x.scanned = append(x.scanned, comps...)
	return x.resultOK, x.resultErr
}

type fakeContext struct {
	flags map[string]string
}

func (f *fakeContext) GetStringFlagValue(flagName string) string {
	return f.flags[flagName]
}
