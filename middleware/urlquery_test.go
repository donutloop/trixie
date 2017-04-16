package middleware_test

import (
	"github.com/donutloop/trixie/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLQuery(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {

		if rv := r.Context().Value(middleware.URLQueryKey); rv != nil {
			if _, ok := rv.(*middleware.Queries); !ok {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

	testHandler := http.HandlerFunc(handler)

	test := httptest.NewServer(middleware.URLQuery()(testHandler))
	defer test.Close()

	response, err := http.Get(test.URL + "?limit=10")
	if err != nil {
		t.Errorf("url query middleware request: %s", err.Error())
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Error("url query middleware request: Unexpected bad request")
	}
}
