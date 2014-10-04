package main

import (
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "github.com/docker/docker/engine"
    "github.com/docker/docker/pkg/version"
    "api"
    "api/server"
)

type Hello struct{}

func (h Hello) ServeHTTP( w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Hello!")

    eng := engine.New()
    var called bool
    eng.Register("info2", func(job *engine.Job) engine.Status {
    	called = true
    	v := &engine.Env{}
    	v.SetInt("Containers", 1)
    	v.SetInt("Images", 300)
    	if _, err := v.WriteTo(job.Stdout); err != nil {
    		return job.Error(err)
    	}
    	return engine.StatusOK
    })

    r2 := serveRequest("GET", "/info2", nil, eng)
    fmt.Fprint(w, r2.Body)
    fmt.Print(r2.Body)
}

func serveRequest(method, target string, body io.Reader, eng *engine.Engine) *httptest.ResponseRecorder {
	return serveRequestUsingVersion(method, target, api.APIVERSION, body, eng)
}

func serveRequestUsingVersion(method, target string, version version.Version, body io.Reader, eng *engine.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRecorder()
	req, err := http.NewRequest(method, target, body)
	if err != nil {
		fmt.Print("error1")
	}
	if err := server.ServeRequest(eng, version, r, req); err != nil {
		fmt.Print("error1")
	}
	return r
}

func main() {
    var h Hello

    host := "localhost:4000"

    http.ListenAndServe(host, h)
    fmt.Print(">>> daemon launched at " + host + "\n")
    fmt.Print(">>>\n")
    fmt.Print(">>> 可測試指令 curl " + host + "\n")
}
