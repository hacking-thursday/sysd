package ui

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"sysd/mods"
)

func Test_ui(t *testing.T) {
	var (
		assert = assert.New(t)
		err    error
		router *mux.Router
		req    *http.Request
		w      *httptest.ResponseRecorder
	)

	//log.SetLevel(log.DebugLevel)

	router, err = mods.CreateRouter(nil)
	assert.NoError(err, "CreateRouter()")

	req, err = mods.NewApiRequest("GET", "/ui", nil)
	assert.NoError(err, "NewApiRequest()")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusMovedPermanently, w.Code, "Response Code")

	req, err = mods.NewApiRequest("GET", "/ui/", nil)
	assert.NoError(err, "NewApiRequest()")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(http.StatusOK, w.Code, "Response Code")

	log.Debug("Test_ui finished")
}
