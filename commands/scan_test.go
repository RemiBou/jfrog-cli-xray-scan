package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleHello(t *testing.T) {
	conf := &scanConfiguration{
		addressee: "World",
		repeat:    1,
	}
	assert.Equal(t, doScan(conf), "Hello World!")
}

func TestComplexHello(t *testing.T) {
	conf := &scanConfiguration{
		addressee:          "World",
		repeat:             3,
		scanCurrentProject: true,
		prefix:             "test: ",
	}
	assert.Equal(t, doScan(conf), "TEST: HELLO WORLD!\nTEST: HELLO WORLD!\nTEST: HELLO WORLD!")
}
