package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sysinfo(t *testing.T) {
	assert := assert.New(t)

	router, err := createRouter()
	assert.NoError(err, "createRouter()")

	req, err := http.NewRequest("GET", *flApiPrefix+"/sysinfo?pretty=1", nil)
	assert.NoError(err, "http.NewRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
	assert.Contains(w.Body.String(), "Uptime", "Sysinfo should contain Uptime")
	assert.Contains(w.Body.String(), "Loads", "Sysinfo should contain Loads")
}
