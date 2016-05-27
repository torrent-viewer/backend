package router

import (
	"fmt"
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
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
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

// AddResource adds a resource to the router
func (router *Router) AddResource(prefix string, resource interface{}) *Router {
	if t, ok := resource.(Listable); ok {
		router.AddRoute(Route{
			Path:    fmt.Sprintf("/%s", prefix),
			Handler: t.RouteList,
			Method:  "GET",
			Name:    fmt.Sprintf("%s.list", prefix),
		})
	}
	if t, ok := resource.(Storable); ok {
		router.AddRoute(Route{
			Path:    fmt.Sprintf("/%s", prefix),
			Handler: t.RouteStore,
			Method:  "POST",
			Name:    fmt.Sprintf("%s.store", prefix),
		})
	}
	if t, ok := resource.(Viewable); ok {
		router.AddRoute(Route{
			Path:    fmt.Sprintf("/%s/{id:[0-9]+}", prefix),
			Handler: t.RouteView,
			Method:  "GET",
			Name:    fmt.Sprintf("%s.view", prefix),
		})
	}
	if t, ok := resource.(Updatable); ok {
		router.AddRoute(Route{
			Path:    fmt.Sprintf("/%s/{id:[0-9]+}", prefix),
			Handler: t.RouteUpdate,
			Method:  "PATCH",
			Name:    fmt.Sprintf("%s.update", prefix),
		})
	}
	if t, ok := resource.(Destroyable); ok {
		router.AddRoute(Route{
			Path:    fmt.Sprintf("/%s/{id:[0-9]+}", prefix),
			Handler: t.RouteDestroy,
			Method:  "DELETE",
			Name:    fmt.Sprintf("%s.delete", prefix),
		})
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
