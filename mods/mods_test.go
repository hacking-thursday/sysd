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
