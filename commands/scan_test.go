package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dummyScanner = func(comps []component) (*ComponentSummaryResult, error) {
	return &ComponentSummaryResult{}, nil
}

func Test_scan(t *testing.T) {
	lines := make(chan string)
	go func() {
		lines <- "gav://org.apache.httpcomponents:httpclient:4.5.9"
		close(lines)
	}()
	err := scan(lines, dummyScanner)

	assert.NoError(t, err)
}
