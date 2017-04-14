package trixie

import (
	"github.com/donutloop/trixie/middleware"
	"net/http"
	"path"
	"strings"
)

// NewRouter returns a new router instance.
func NewRouter() *Router {
	return &Router{}
}

// Router registers routes to be matched and dispatches a handler.
//
// It implements the http.Handler interface, so it can be registered to serve
// requests:
//
//     var router = mux.NewRouter()
//
//     func main() {
//         http.Handle("/", router)
//     }
//
// This will send all incoming requests to the router.
type Router struct {
	// Configurable Handler to be used when no route matches.
	NotFoundHandler http.Handler

	// This defines the flag for new routes.
	StrictSlash bool
	// This defines the flag for new routes.
	SkipClean bool
	// This defines a flag for all routes.
	UseEncodedPath bool
	// This defines a flag for all routes.
	CaseSensitiveURL bool
	// this builds a tree
	treeConstructor func() RouteTreeInterface
	// This defines the tree for routes.
	tree RouteTreeInterface
	// this builds a route
	routeConstructor func() RouteInterface

	// The middleware stack
	middlewares []middleware.Middleware
}

// Use appends a middleware handler to the Mux middleware stack.
func (router *Router) Use(middlewares ...middleware.Middleware) {
	router.middlewares = append(router.middlewares, middlewares...)
}

// UseRoute that you can use diffrent route versions
// See RouteInterface for more details (route.go)
func (r *Router) UseRoute(constructer func() RouteInterface) {
	r.routeConstructor = constructer
}

// UseTree that you can use diffrent tree versions
// See TreeInterface for more details (tree.go)
func (r *Router) UseTree(constructer func() RouteTreeInterface) {
	r.treeConstructor = constructer
}

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.GetVars(req).Get(":number") or mux.GetVars(req).GetAll()
//
// and the route queires can be retrieved calling
// mux.GetQueries(req).Get(":number") or mux.GetQueries(req).GetAll()
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	method := Methods.lookup(req.Method)

	if method == methodNotFound {
		r.notFoundHandler().ServeHTTP(w, req)
		return
	}

	if !r.SkipClean {

		p := req.URL.Path

		if r.UseEncodedPath {
			p = req.URL.EscapedPath()
		}

		// Clean path to canonical form and redirect.
		if cp := cleanPath(p); p != p {
			w.Header().Set("Location", cp)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	if !r.CaseSensitiveURL {
		req.URL.Path = strings.ToLower(req.URL.Path)
	}

	route, params, err := r.tree.Find(r.tree.GetRoot(), req.URL.Path)
	if err != nil || !route.HasHandler(method) {
		r.notFoundHandler().ServeHTTP(w, req)
		return
	}

	req = AddCurrentRoute(req, route)
	req = AddRouteParameters(req, params)

	middleware.Stack(r.middlewares...).Then(route.GetHandler(method)).ServeHTTP(w, req)
}

func (r *Router) notFoundHandler() http.Handler {
	if r.NotFoundHandler == nil {
		return http.NotFoundHandler()
	}

	return r.NotFoundHandler
}

// cleanPath returns the canonical path for p, eliminating . and .. elements.
// Borrowed from the net/http package.
// /net/http/server.go
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

// RegisterRoute registers and validates a new route
func (r *Router) RegisterRoute(route RouteInterface) {

	if r.tree == nil {
		r.tree = r.treeConstructor()
	}

	r.tree.Insert(route)
}

func (r *Router) ValidateRoute(route RouteInterface) {
	for _, validator := range Validatoren {
		err := validator.Validate(route)

		if err != nil {
			panic(err.Error())
		}
	}
}

// Handle registers a new route with a matcher for the URL path.
func (r *Router) Handle(method string, pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(method, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// HandleFunc registers a new route with a matcher for the URL path.
func (r *Router) HandleFunc(method string, pattern string, handler http.Handler) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandler(method, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Get registers a new get route for the URL path
func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodGet, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Put registers a new put route for the URL path
func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodPut, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Post registers a new post route for the URL path
func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodPost, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Delete registers a new delete route for the URL path
func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodDelete, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Patch registers a new patch route for the URL path
func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodPatch, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Options registers a new options route for the URL path
func (r *Router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodOptions, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

// Head registers a new head route for the URL path
func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	route.AddHandlerFunc(http.MethodHead, handler)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

func (r *Router) Path(pattern string, callback func(route RouteInterface)) RouteInterface {
	route := r.routeConstructor()
	route.SetPattern(pattern)
	callback(route)
	r.ValidateRoute(route)
	r.RegisterRoute(route)
	return route
}

var Methods = newMethods()

type methods struct {
	ms map[string]Method
}

func newMethods() *methods {
	return &methods{
		ms: methodsMap,
	}
}

// lookup check if method exists when return Method else return MethodNotFound
func (m *methods) lookup(method string) Method {

	if value, found := m.ms[method]; found {
		return value
	}

	return methodNotFound
}

func (m *methods) lookupID(method Method) string {
	for k, v := range m.ms {
		if method == v {
			return k
		}
	}
	return ""
}

func (m *methods) Set(method string, methodN Method) {
	m.ms[method] = methodN
}

func (m *methods) Delete(method string) {
	delete(m.ms, method)
}

// Methods a map of all standard methods
var methodsMap = map[string]Method{
	http.MethodGet:     MethodGet,
	http.MethodPost:    MethodPost,
	http.MethodPut:     MethodPut,
	http.MethodDelete:  MethodDelete,
	http.MethodPatch:   MethodPatch,
	http.MethodOptions: MethodOptions,
	http.MethodHead:    MethodHead,
}

var methodNotFound Method = -1
