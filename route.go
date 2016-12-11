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
)

type RouteInterface interface {
	Kind() uint8
	AddHandler(method Method, handler http.Handler) RouteInterface
	SetPattern(string) RouteInterface
	GetHandler(Method) http.Handler
	HasHandler(Method) bool
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

func (r *Route) Kind() uint8 {
	return 1
}

func (r *Route) AddHandler(method Method, handler http.Handler) RouteInterface {
	r.handlers[method] = handler
	return r
}

func (r *Route) GetHandler(method Method) http.Handler {
	return r.handlers[method]
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
