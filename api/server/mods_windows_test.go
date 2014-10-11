package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_osver(t *testing.T) {
	assert := assert.New(t)

	router, err := createRouter()
	assert.NoError(err, "createRouter()")

	req, err := http.NewRequest("GET", *flApiPrefix+"/osver?pretty=1", nil)
	assert.NoError(err, "http.NewRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
	assert.Contains(w.Body.String(), "windows", "Type should be windows")
}
