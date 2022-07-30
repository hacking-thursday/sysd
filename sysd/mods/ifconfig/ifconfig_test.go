package ifconfig

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"sysd/mods"
)

func Test_ifconfig(t *testing.T) {
	assert := assert.New(t)

	router, err := mods.CreateRouter(nil)
	assert.NoError(err, "CreateRouter()")

	req, err := mods.NewApiRequest("GET", "/ifconfig?pretty=1", nil)
	assert.NoError(err, "NewApiRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
}
