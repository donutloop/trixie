package tmux

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeTest struct {
	title      string
	path       string
	method     string
	statusCode int
	kind       string
	queries    map[string][]string
	vars       map[string]string
	route      func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request))
}

func TestPath(t *testing.T) {

	tests := []routeTest{
		{
			title:      "(GET) Path route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			route: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Get(path, handler)
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %s path: %s method %s kind: %s", test.title, test.path, test.method, test.kind), func(t *testing.T) {
			code, message, ok := testRoute(test)

			if !ok {
				t.Errorf("Expected status code %v, Actucal status code %v, Actucal message %v", test.statusCode, code, message)
			}
		})
	}
}

func testRoute(rt routeTest) (int, string, bool) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("succesfully"))
	}

	r := Classic()
	rt.route(r, rt.path, rt.method, handler)

	req, _ := http.NewRequest(rt.method, "http://localhost"+rt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, "", false
	}

	if res.Code != rt.statusCode || content.String() != "succesfully" {
		return res.Code, content.String(), false
	}

	return res.Code, content.String(), true
}
