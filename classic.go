package tmux

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewTree(NewNode))
	router.UseRoute(NewRoute)
	return router
}
