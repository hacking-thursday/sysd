package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hacking-thursday/sysd/mods"
)

func Test_ping(t *testing.T) {
	assert := assert.New(t)

	router, err := mods.CreateRouter(nil)
	assert.NoError(err, "CreateRouter()")

	req, err := mods.NewApiRequest("GET", "/ping", nil)
	assert.NoError(err, "NewApiRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
	assert.Equal("pong", w.Body.String(), "Response Body")
}
