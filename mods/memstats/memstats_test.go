package memstats

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hacking-thursday/sysd/mods"
)

func Test_memstats(t *testing.T) {
	assert := assert.New(t)

	router, err := mods.CreateRouter(nil)
	assert.NoError(err, "CreateRouter()")

	req, err := mods.NewApiRequest("GET", "/memstats?pretty=1", nil)
	assert.NoError(err, "http.NewRequest()")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")
}
