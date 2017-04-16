package trixie

import (
	"context"
	"github.com/donutloop/trixie/middleware"
	"net/http"
)

// Context keys
const (
	queriesKey middleware.ContextKey = "urlqueryKey"
	routeKey                         = "routeKey"
	paramKey                         = "paramKey"
)

// GetQueries returns the query variables for the current request.
func GetQueries(r *http.Request) *middleware.Queries {
	if value := r.Context().Value(queriesKey); value != nil {
		return value.(*middleware.Queries)
	}

	return nil
}

// CurrentRoute returns the matched route for the current request.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns
func GetCurrentRoute(r *http.Request) RouteInterface {
	if rv := r.Context().Value(routeKey); rv != nil {
		return rv.(RouteInterface)
	}

	return nil
}

// AddCurrentRoute adds a route instance to the current request context
func AddCurrentRoute(r *http.Request, route RouteInterface) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), routeKey, route))
}

// AddCurrentRoute adds parameters of path to the current request context
func AddRouteParameters(r *http.Request, params map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), paramKey, params))
}

// GetRouteParameter returns the parameters of route for a given request
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns
func GetRouteParameters(r *http.Request) map[string]string {
	if rv := r.Context().Value(paramKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}
