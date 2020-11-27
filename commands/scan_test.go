package commands

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIfComponentFlagScanIt(t *testing.T) {
	client := &fakeXrayClient{
		resultOK: &ComponentSummaryResult{},
	}
	context := &fakeContext{flags: map[string]string{componentFlagKey: "acomponent"}}
	s := xrayScanner{xrayClient: client, stdin: nil}
	s.scanCmd(context)
	require.ElementsMatch(t, []component{"acomponent"}, client.lastScan)
}

func TestIfStdInMVNDependencyListScanAllUntilEOF(t *testing.T) {
	client := &fakeXrayClient{resultOK: &ComponentSummaryResult{}}
	context := &fakeContext{flags: map[string]string{}}
	stdin := bytes.NewBufferString(
		"[INFO]    org.slf4j:slf4j-ext:jar:1.7.26:compile -- module slf4j.ext (auto)\n" +
			"[INFO]    org.jolokia:jolokia-core:jar:1.6.2:compile -- module jolokia.core (auto)\n" +
			"[INFO]    com.googlecode.json-simple:json-simple:jar:1.1.1:compile -- module json.simple (auto)\n" +
			"[INFO]    io.jaegertracing:jaeger-client:jar:1.2.0:compile -- module jaeger.client (auto)\n" +
			"[INFO]    io.jaegertracing:jaeger-thrift:jar:1.2.0:compile -- module jaeger.thrift (auto)\n" +
			"[INFO]    org.apache.thrift:libthrift:jar:0.13.0:compile -- module libthrift (auto)\n" +
			"[INFO]    com.google.code.findbugs:jsr305:jar:2.0.0:compile -- module jsr305 (auto)\n" +
			"[INFO] \n" +
			"[INFO] \n" +
			"[INFO] ----------------< org.jfrog.access:access-coverage-all >----------------\n" +
			"[INFO] Building JFrog Access OSS Coverage 7.x-SNAPSHOT                  [21/21]\n" +
			"[INFO] --------------------------------[ jar ]---------------------------------\n" +
			"[INFO] \n")
	s := xrayScanner{xrayClient: client, stdin: stdin}
	s.scanCmd(context)
	require.ElementsMatch(t, []component{
		"gav:org.jolokia:jolokia-core:1.6.2",
		"gav:org.slf4j:slf4j-ext:1.7.26",
		"gav:com.googlecode.json-simple:json-simple:1.1.1",
		"gav:io.jaegertracing:jaeger-client:1.2.0",
		"gav:io.jaegertracing:jaeger-thrift:1.2.0",
		"gav:com.google.code.findbugs:jsr305:2.0.0",
		"gav:org.apache.thrift:libthrift:0.13.0"}, client.lastScan)

}

func TestIfStdInGoDependencyListScanAllUntilEOF(t *testing.T) {

}

func TestValidateOutput(t *testing.T) {

}

func TestOnScanError(t *testing.T) {

}

type fakeXrayClient struct {
	lastScan  []component
	resultOK  *ComponentSummaryResult
	resultErr error
}

func (x *fakeXrayClient) scan(comps []component) (*ComponentSummaryResult, error) {
	x.lastScan = comps
	return x.resultOK, x.resultErr
}

type fakeContext struct {
	flags map[string]string
}

func (f *fakeContext) GetStringFlagValue(flagName string) string {
	return f.flags[flagName]
}
