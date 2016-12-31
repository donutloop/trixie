package tmux

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewRadixTree(NewNode))
	router.UseRoute(NewRoute)
	return router
}
