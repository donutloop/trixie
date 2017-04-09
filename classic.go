package tmux

import "github.com/donutloop/tmux/middleware"

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewTree(NewNode))
	router.UseRoute(NewRoute)
	router.Use(middleware.URLQuery())
	return router
}
