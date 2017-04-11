package tmux

import (
	"context"
	"github.com/donutloop/tmux/middleware"
	"net/http"
)

const (
	queriesKey middleware.ContextKey = "urlqueryKey"
	routeKey                         = "routeKey"
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

func AddCurrentRoute(r *http.Request, route RouteInterface) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), routeKey, route))
}

func AddRouteParameters(r *http.Request, params map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), routeKey, params))
}
