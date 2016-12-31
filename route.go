package tmux

import "net/http"

type Method int8

type Handlers map[Method]http.Handler

// Method of request
const (
	MethodGet Method = iota
	MethodPost
	MethodPut
	MethodDelete
	MethodOptions
	MethodPatch
	MethodHead
)

type RouteInterface interface {
	AddHandler(method string, handler http.Handler) RouteInterface
	AddHandlerFunc(method string, handler func(http.ResponseWriter, *http.Request)) RouteInterface
	SetPattern(string) RouteInterface
	GetPattern() string
	GetHandler(Method) http.Handler
	HasHandler(Method) bool
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
	r.handlers[methods.lookup(method)] = http.HandlerFunc(handler)
	return r
}

func (r *Route) AddHandler(method string, handler http.Handler) RouteInterface {
	r.handlers[methods.lookup(method)] = handler
	return r
}

func (r *Route) GetHandler(method Method) http.Handler {
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

func (r *Route) HasHandler(method Method) bool {
	if _, found := r.handlers[method]; found {
		return true
	}
	return false
}
