package server2

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_createRouter(t *testing.T) {
	var (
		err error
		r   *mux.Router
	)
	assert := assert.New(t)

	r, err = createRouter()
	assert.NotNil(r, "createRouter")
	assert.NoError(err, "createRouter")
}

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
