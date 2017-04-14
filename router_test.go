package trixie

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeTestCase struct {
	title        string
	path         string
	rawPath      string
	method       string
	countOfParam int
	statusCode   int
	defineRoute  func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request))
}

func TestPath(t *testing.T) {
	tests := []routeTestCase{
		{
			title:      "(GET) route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			defineRoute: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Get(path, handler)
			},
		},
		{
			title:      "(Path) route with single path",
			path:       "/api/",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			defineRoute: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				r.Path(path, func(route RouteInterface) {
					route.AddHandlerFunc(method, handler)
				})
			},
		},
		{
			title:      "Path (Method not found)",
			path:       "/api/",
			method:     "GETT",
			statusCode: http.StatusNotFound,
			defineRoute: func(r *Router, path string, method string, handler func(w http.ResponseWriter, r *http.Request)) {
				Methods.Set(method)
				r.Path(path, func(route RouteInterface) {
					route.AddHandlerFunc(method, handler)
				})
				Methods.Delete(method)
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test: %s path: %s method %s", test.title, test.path, test.method), func(t *testing.T) {
			statusCode, message := testSingleRoute(test)

			switch test.statusCode {
			case http.StatusNotFound:
				if statusCode != test.statusCode {
					t.Errorf("Expected status code %v, Actucal status code %v", test.statusCode, statusCode)
				}
			default:
				if statusCode != test.statusCode || message != "succesfully" {
					t.Errorf("Expected status code %v, Actucal status code %v", test.statusCode, statusCode)
				}
			}
		})
	}
}

func testSingleRoute(rt routeTestCase) (int, string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("succesfully"))
	}

	r := Classic()
	rt.defineRoute(r, rt.path, rt.method, handler)

	req, _ := http.NewRequest(rt.method, "http://localhost"+rt.path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var content bytes.Buffer
	_, err := io.Copy(&content, res.Body)

	if err != nil {
		return -1, "Error while reading of respone body"
	}

	return res.Code, content.String()
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
		"/api/article/:number/comment/9": {
			key:  "8",
			path: "/api/article/9/comment/9",
		},
		"/api/article/:number/questions/:number": {
			key:  "9",
			path: "/api/article/97/questions/4",
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
