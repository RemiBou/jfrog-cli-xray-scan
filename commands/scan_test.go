package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var dummyScanner = func(comps []component) (*ComponentSummaryResult, error) {
	// TODO: instead deserialize json files from FS
	return &ComponentSummaryResult{}, nil
}

func Test_scan(t *testing.T) {
	lines := make(chan string)
	go func() {
		lines <- "gav://org.apache.httpcomponents:httpclient:4.5.9"
		lines <- "gav://org.codehaus.plexus:plexus-utils:3.2.1"
		close(lines)
	}()
	err := scan(lines, dummyScanner)

	assert.NoError(t, err)
}
