package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Router is a resource-oriented router
type Router struct {
	mux         *mux.Router
	middlewares []Middleware
}

// Route is a URL path with some context added
type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  string
	Name    string
}

// Routes is an array of Route
type Routes []Route

type Middleware func(http.Handler) http.Handler

// NewRouter creates a new Router instance
func NewRouter() *Router {
	router := new(Router)
	router.mux = mux.NewRouter().StrictSlash(true)
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var handler http.Handler = router.mux
	for _, mw := range router.middlewares {
		handler = mw(handler)
	}
	handler.ServeHTTP(w, r)
}

// AddRoute adds a route to the router
func (router *Router) AddRoute(route Route) *Router {
	log.Printf("Registering route %-25s: %10s %s\n", route.Name, route.Method, route.Path)
	router.mux.
		Path(route.Path).
		HandlerFunc(route.Handler).
		Methods(route.Method).
		Name(route.Name)
	return router
}

// AddRoutes adds an slice of routes to the router
func (router *Router) AddRoutes(routes Routes) *Router {
	for _, route := range routes {
		router.AddRoute(route)
	}
	return router
}

// Use adds a middleware to the router
func (router *Router) Use(mw Middleware) *Router {
	router.middlewares = append(router.middlewares, mw)
	return router
}

// Vars returns the variables found in the current route URL
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
