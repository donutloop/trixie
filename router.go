package tmux

import (
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
}

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

	method := methods.lookupMethod(req.Method)

	if method == methodNotFound {
		r.notFoundHandler().ServeHTTP(w, req)
	}

	if !r.SkipClean {

		path := req.URL.Path

		if r.UseEncodedPath {
			path = req.URL.EscapedPath()
		}

		// Clean path to canonical form and redirect.
		if p := cleanPath(path); p != path {
			w.Header().Set("Location", p)
			w.WriteHeader(http.StatusMovedPermanently)
			return
		}
	}

	if !r.CaseSensitiveURL {
		req.URL.Path = strings.ToLower(req.URL.Path)
	}

	route := r.tree.Find(r.tree.GetRoot(), method, req.URL.Path)

	if route == nil {
		r.notFoundHandler().ServeHTTP(w, req)
		return
	}

	RequestContext.AddCurrentRoute(req, route)
	RequestContext.AddQueries(req)

	route.GetHandler(method).ServeHTTP(w, req)
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
func (r *Router) RegisterRoute(method Method, pattern string, handler http.Handler) RouteInterface {

	if r.tree == nil {
		r.tree = r.treeConstructor()
	}

	return r.tree.Insert(method, pattern, handler)
}

// Handle registers a new route with a matcher for the URL path.
func (r *Router) Handle(method Method, pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(method, pattern, http.HandlerFunc(handler))
}

// Get registers a new get route for the URL path
func (r *Router) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodGet, pattern, http.HandlerFunc(handler))
}

// Put registers a new put route for the URL path
func (r *Router) Put(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodPut, pattern, http.HandlerFunc(handler))
}

// Post registers a new post route for the URL path
func (r *Router) Post(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodPost, pattern, http.HandlerFunc(handler))
}

// Delete registers a new delete route for the URL path
func (r *Router) Delete(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodDelete, pattern, http.HandlerFunc(handler))
}

// Patch registers a new patch route for the URL path
func (r *Router) Patch(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodPatch, pattern, http.HandlerFunc(handler))
}

// Options registers a new options route for the URL path
func (r *Router) Options(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodOptions, pattern, http.HandlerFunc(handler))
}

// Head registers a new head route for the URL path
func (r *Router) Head(pattern string, handler func(http.ResponseWriter, *http.Request)) RouteInterface {
	return r.RegisterRoute(MethodHead, pattern, http.HandlerFunc(handler))
}

var methods = newMethods()

type Methods struct {
	ms map[string]Method
}

func newMethods() *Methods {
	return &Methods{
		ms: methodsMap,
	}
}

// lookupMethod check if method exists when return Method else return MethodNotFound
func (m *Methods) lookupMethod(method string) Method {

	if method, found := m.ms[method]; found {
		return method
	}

	return methodNotFound
}

func (m *Methods) Set(method string, methodN Method) {
	m.ms[method] = methodN
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
