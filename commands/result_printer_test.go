package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func Test_resultPrinter_print(t *testing.T) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", testData, "few_components.json"))

	bufferString := bytes.NewBufferString("")
	printer, err := newPrinter(bufferString)
	require.NoError(t, err)

	result := &ComponentSummaryResult{}
	err = json.Unmarshal(file, result)
	require.NoError(t, err)

	err = printer.print(*result)
	require.NoError(t, err)
	expected, err := ioutil.ReadFile(fmt.Sprintf("%voutput/%v", testData, "few_components.txt"))
	require.Equal(t, string(expected), bufferString.String())
}
