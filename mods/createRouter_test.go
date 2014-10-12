package mods

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

	r, err = CreateRouter(nil)
	assert.NotNil(r, "CreateRouter")
	assert.NoError(err, "CreateRouter")
}
