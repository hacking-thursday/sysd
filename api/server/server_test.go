package server

import (
	"net/http"
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

func Test_marshal(t *testing.T) {
	assert := assert.New(t)

	req, err := http.NewRequest("GET", *flApiPrefix+"/memstats?pretty=1", nil)
	assert.NoError(err, "http.NewRequest()")

	b, err := marshal(req, req)
	assert.Contains(string(b), "\t", "marshal pretty should contains <TAB>")

	req, err = http.NewRequest("GET", *flApiPrefix+"/memstats", nil)
	assert.NoError(err, "http.NewRequest()")

	b, err = marshal(req, req)
	assert.NotContains(string(b), "\t", "marshal should not contains <TAB>")
}
