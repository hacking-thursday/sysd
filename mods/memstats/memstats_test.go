package memstats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_memstats(t *testing.T) {
	assert := assert.New(t)

	router, err := createRouter()
	assert.NoError(err, "createRouter()")

	req, err := http.NewRequest("GET", *flApiPrefix+"/memstats?pretty=1", nil)
	assert.NoError(err, "http.NewRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
}
