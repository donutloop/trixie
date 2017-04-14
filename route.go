package trixie

import "net/http"

type method string

type Handlers map[string]http.Handler

type RouteInterface interface {
	AddHandler(string, http.Handler) RouteInterface
	AddHandlerFunc(string, func(http.ResponseWriter, *http.Request)) RouteInterface
	SetPattern(string) RouteInterface
	GetPattern() string
	GetHandler(string) http.Handler
	HasHandler(string) bool
	GetHandlers() Handlers
	AddHandlers(Handlers) RouteInterface
}

func NewRoute() RouteInterface {
	return &Route{
		handlers: Handlers{},
	}
}

type Route struct {
	handlers Handlers
	pattern  string
}

func (r *Route) AddHandlerFunc(method string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	r.handlers[method] = http.HandlerFunc(handler)
	return r
}

func (r *Route) AddHandler(method string, handler http.Handler) RouteInterface {
	r.handlers[method] = handler
	return r
}

func (r *Route) GetHandler(method string) http.Handler {
	return r.handlers[method]
}

func (r *Route) AddHandlers(handlers Handlers) RouteInterface {

	for method, handler := range handlers {
		r.handlers[method] = handler
	}

	return r
}

func (r *Route) GetHandlers() Handlers {
	return r.handlers
}

func (r *Route) SetPattern(pattern string) RouteInterface {
	r.pattern = pattern
	return r
}

func (r *Route) GetPattern() string {
	return r.pattern
}

func (r *Route) HasHandler(method string) bool {
	if _, found := r.handlers[method]; found {
		return true
	}
	return false
}
