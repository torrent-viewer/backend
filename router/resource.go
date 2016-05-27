package router

import (
	"net/http"
)

type Listable interface {
	RouteList(w http.ResponseWriter, r *http.Request)
}

type Storable interface {
	RouteStore(w http.ResponseWriter, r *http.Request)
}

type Viewable interface {
	RouteView(w http.ResponseWriter, r *http.Request)
}

type Updatable interface {
	RouteUpdate(w http.ResponseWriter, r *http.Request)
}

type Destroyable interface {
	RouteDestroy(w http.ResponseWriter, r *http.Request)
}