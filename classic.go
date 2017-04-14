package trixie

import "github.com/donutloop/trixie/middleware"

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewTree(NewNode))
	router.UseRoute(NewRoute)
	router.Use(middleware.URLQuery())
	return router
}
