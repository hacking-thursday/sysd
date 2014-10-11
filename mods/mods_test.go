package mods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	var err error
	assert := assert.New(t)

	err = Register("GET", "/TEST", nil)
	assert.NoError(err, "Register")

	err = Register("GET", "/TEST", nil)
	assert.Error(err, "Register twice")
}

func Test_Marshal(t *testing.T) {
	assert := assert.New(t)

	req, err := NewApiRequest("GET", "/memstats?pretty=1", nil)
	assert.NoError(err, "NewApiRequest()")

	b, err := Marshal(req, req)
	assert.Contains(string(b), "\t", "marshal pretty should contains <TAB>")

	req, err = NewApiRequest("GET", "/memstats", nil)
	assert.NoError(err, "NewApiRequest()")

	b, err = Marshal(req, req)
	assert.NotContains(string(b), "\t", "marshal should not contains <TAB>")
}
