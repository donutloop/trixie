package tmux

import (
	"context"
	"net/http"
	"github.com/donutloop/tmux/middleware"
)


const (
	queriesKey middleware.ContextKey = "urlqueryKey"
	routeKey = "routeKey"
)

type ReqContext struct{}

// GetQueries returns the query variables for the current request.
func (c ReqContext) GetQueries(r *http.Request) *middleware.Queries {
	if value := r.Context().Value(queriesKey); value != nil {
		return value.(*middleware.Queries)
	}

	return nil
}

// CurrentRoute returns the matched route for the current request.
// This only works when called inside the handler of the matched route
// because the matched route is stored in the request context which is cleared
// after the handler returns
func (c ReqContext) GetCurrentRoute(r *http.Request) RouteInterface {
	if rv := r.Context().Value(routeKey); rv != nil {
		return rv.(RouteInterface)
	}

	return nil
}

func (c ReqContext) AddCurrentRoute(r *http.Request, route RouteInterface) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), routeKey, route))
}
