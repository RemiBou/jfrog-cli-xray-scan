package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

const testData = "testData/"

func Test_resultPrinter_print_no_issues(t *testing.T) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, "xray_response.json"))

	bufferString := bytes.NewBufferString("")
	printer, err := newPrinter(bufferString, printerConfig{printNoIssues: true})
	require.NoError(t, err)

	result := &ComponentSummaryResult{}
	err = json.Unmarshal(file, result)
	require.NoError(t, err)

	err = printer.print(*result)
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(fmt.Sprintf("%voutput/%v", testData, "few_components_all_issues.txt"))
	require.Equal(t, string(expected), bufferString.String())
}

func Test_resultPrinter_print(t *testing.T) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, "xray_response.json"))

	bufferString := bytes.NewBufferString("")
	printer, err := newPrinter(bufferString, printerConfig{printNoIssues: false})
	require.NoError(t, err)

	result := &ComponentSummaryResult{}
	err = json.Unmarshal(file, result)
	require.NoError(t, err)

	err = printer.print(*result)
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(fmt.Sprintf("%voutput/%v", testData, "few_components.txt"))
	require.Equal(t, string(expected), bufferString.String())
}
