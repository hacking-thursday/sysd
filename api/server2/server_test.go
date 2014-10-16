package server2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseAddr(t *testing.T) {
	var (
		err   error
		proto string
		addr  string
	)
	assert := assert.New(t)

	proto, addr, err = parseAddr("tcp://0.0.0.0:8000")
	assert.NoError(err, "parseAddr")
	assert.Equal("tcp", proto, "parseAddr")
	assert.Equal("0.0.0.0:8000", addr, "parseAddr")

	proto, addr, err = parseAddr("0.0.0.0:8000")
	assert.Error(err, "parseAddr")

	proto, addr, err = parseAddr("unix://0.0.0.0:8000")
	assert.Error(err, "parseAddr")
}
