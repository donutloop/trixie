package tmux

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

type contextKey int

const (
	queriesKey contextKey = iota
	routeKey
)

type requestContext struct{}

// GetQueries returns the query variables for the current request.
func (c requestContext) GetQueries(r *http.Request) queries {
	if value := r.Context().Value(queriesKey); value != nil {
		return value.(queries)
	}

	return nil
}

// CurrentRoute returns the matched route for the current request.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns
func (c requestContext) GetCurrentRoute(r *http.Request) RouteInterface {
	if rv := r.Context().Value(routeKey); rv != nil {
		return rv.(RouteInterface)
	}

	return nil
}

func (c requestContext) AddQueries(r *http.Request) *http.Request {
	queries, err := extractQueries(r)

	if err != nil || 0 == queries.Count() {
		return r
	}

	return r.WithContext(context.WithValue(r.Context(), queriesKey, queries))
}

func (c requestContext) AddCurrentRoute(r *http.Request, route RouteInterface) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), routeKey, route))
}

type queries map[string][]string

// Get return the key value, of the current *http.Request queries
func (q queries) Get(key string) []string {
	if value, found := q[key]; found {
		return value
	}
	return []string{}
}

// Get returns all queries of the current *http.Request queries
func (q queries) GetAll() map[string][]string {
	return q
}

// Count returns count of the current *http.Request queries
func (q queries) Count() int {
	return len(q)
}

// extractQueries extract queries of the given *http.Request
func extractQueries(req *http.Request) (queries, error) {

	queriesRaw, err := url.ParseQuery(req.URL.RawQuery)

	if err != nil {
		return nil, err
	}

	queries := queries(map[string][]string{})

	if 0 == len(queriesRaw) {
		return queries, nil
	}

	for k, v := range queriesRaw {
		for _, item := range v {
			values := strings.Split(item, ",")
			queries[k] = append(queries[k], values...)
		}
	}

	return queries, nil
}

var RequestContext = &requestContext{}
