package tmux

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewRadixTree(NewNode, NewRoute))
	return router
}
