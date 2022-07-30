package apilist

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"sysd/mods"
)

func Test_apilist(t *testing.T) {
	assert := assert.New(t)

	router, err := mods.CreateRouter(nil)
	assert.NoError(err, "CreateRouter()")

	req, err := mods.NewApiRequest("GET", "/apilist?pretty=1", nil)
	assert.NoError(err, "NewApiRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
	assert.Contains(w.Body.String(), "/apilist", "Response should contain apilist")
}
