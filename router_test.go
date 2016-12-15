package tmux

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeTestCase struct {
	title      string
	path       string
	method     string
	statusCode int
	queries    map[string][]string
	vars       map[string]string
	route      func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request))
}

func TestPath(t *testing.T) {

	tests := []routeTestCase{
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
		t.Run(fmt.Sprintf("Test: %s path: %s method %s", test.title, test.path, test.method), func(t *testing.T) {
			code, message, ok := testSingleRoute(test)

			if !ok {
				t.Errorf("Expected status code %v, Actucal status code %v, Actucal message %v", test.statusCode, code, message)
			}
		})
	}
}

func testSingleRoute(rt routeTestCase) (int, string, bool) {
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

func TestRouterWithMultiRoutes(t *testing.T) {
	router := Classic()

	handler := func(key string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(key))
		}
	}

	paths := map[string]struct {
		key  string
		path string
	}{
		"/api/user/1/comment/1": {
			key:  "0",
			path: "/api/user/1/comment/1",
		},
		"/api/echo": {
			key:  "1",
			path: "/api/echo",
		},
		"/api/user/:number/comments": {
			key:  "2",
			path: "/api/user/2/comments",
		},
		"/api/user/donutloop": {
			key:  "3",
			path: "/api/user/donutloop",
		},
		"/api/user/6": {
			key:  "4",
			path: "/api/user/6",
		},
		"/api/article/golang": {
			key:  "6",
			path: "/api/article/golang",
		},
		"/api/article/7": {
			key:  "7",
			path: "/api/article/7",
		},
		"/api/article/97/comment/9": {
			key:  "8",
			path: "/api/article/97/comment/9",
		},
	}

	for path, pathInfo := range paths {
		router.Get(path, handler(pathInfo.key))
	}

	server := httptest.NewServer(router)

	for rawPath, pathInfo := range paths {
		url := server.URL + pathInfo.path
		t.Run(fmt.Sprintf("RawPath: %s, Path: %s Url: %s", rawPath, pathInfo.path, url), func(t *testing.T) {
			res, err := http.Get(url)

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			var content bytes.Buffer
			_, err = io.Copy(&content, res.Body)

			if err != nil {
				t.Errorf("Unexpected error (%s)", err.Error())
			}

			defer res.Body.Close()

			if content.String() != pathInfo.key {
				t.Errorf("Path %s: Unexpected path key (Expected path key: %s, Actucal path key %s)", rawPath, pathInfo.key, content.String())
			}
		})
	}
}
