package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hacking-thursday/sysd/mods"
)

func Test_osver(t *testing.T) {
	assert := assert.New(t)

	router, err := mods.CreateRouter()
	assert.NoError(err, "CreateRouter()")

	req, err := mods.NewApiRequest("GET", "/osver?pretty=1", nil)
	assert.NoError(err, "NewApiRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
	assert.Contains(w.Body.String(), "windows", "Type should be windows")
}
