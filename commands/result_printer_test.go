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

func Test_resultPrinter_print(t *testing.T) {
	tests := []struct {
		name          string
		printNoIssues bool
		expected      string
	}{
		{name: "Print components including with no issues", printNoIssues: true,
			expected: "few_components_all_issues.txt"},
		{name: "Don't print components with no issues", printNoIssues: false,
			expected: "few_components.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, "xray_response.json"))

			bufferString := bytes.NewBufferString("")
			printer := newPrinter(bufferString, printerConfig{printNoIssues: tt.printNoIssues})
			require.NoError(t, err)

			result := &ComponentSummaryResult{}
			err = json.Unmarshal(file, result)
			require.NoError(t, err)

			printer.print(*result)
			printer.flush()

			expected, err := ioutil.ReadFile(fmt.Sprintf("%voutput/%v", testData, tt.expected))
			require.Equal(t, string(expected), bufferString.String())
		})
	}
}
