package tmux

func Classic() *Router {
	router := NewRouter()
	router.UseTree(NewRadixTree(NewNode))
	return router
}
