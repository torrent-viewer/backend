package show

import (
	"fmt"
	"net/http"

	"github.com/torrent-viewer/backend/datastore"
	"github.com/torrent-viewer/backend/responses"
	"github.com/torrent-viewer/backend/requests"
	"github.com/torrent-viewer/backend/herr"
)

// ShowsList is the HTTP endpoint used to create list Shows instances
func (ShowResource) RouteList(w http.ResponseWriter, r *http.Request) {
	var entries Shows
	if err := datastore.FetchEntities(&entries); err != nil {
		responses.SendError(w, *err)
		return
	}
	serialized := make([]interface{}, len(entries), len(entries))
	for i, e := range entries {
		serialized[i] = e
	}
	responses.SendEntities(w, serialized)
}

// ShowsStore is the HTTP endpoint used to create new Shows instances
func (ShowResource) RouteStore(w http.ResponseWriter, r *http.Request) {
	var show Show
	if err := requests.ReceiveEntity(r, &show); err != nil {
		responses.SendError(w, *err)
		return
	}
	if err := datastore.StoreEntity(&show); err != nil {
		responses.SendError(w, *err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/shows/%d", show.ID))
	responses.SendEntity(w, &show, http.StatusCreated)
}

// ShowsView is the HTTP endpoint used to show Shows instance by ID
func (ShowResource) RouteView(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		responses.SendError(w, *err)
		return
	}
	var show Show
	if err := datastore.FetchEntity(&show, id); err != nil {
		responses.SendError(w, *err)
		return
	}
	responses.SendEntity(w, &show, http.StatusOK)
}

// ShowsUpdate is the HTTP endpoint used to update a Show instance by its ID
func (ShowResource) RouteUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		responses.SendError(w, *err)
		return
	}
	var show Show
	if err := datastore.FetchEntity(&show, id); err != nil {
		responses.SendError(w, *err)
		return
	}
	if err := requests.ReceiveEntity(r, &show); err != nil {
		responses.SendError(w, *err)
		return
	}
	if show.ID != id {
		responses.SendError(w, herr.UnmatchingIDsError)
		return
	}
	if err := datastore.UpdateEntity(&show); err != nil {
		responses.SendError(w, *err)
		return
	}
	responses.SendNoContent(w)
}

// ShowsDestroy is the HTTP endpoint used to delete a Show instance by its ID
func (ShowResource) RouteDestroy(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		responses.SendError(w, *err)
		return
	}
	show := Show{
		ID: id,
	}
	if err := datastore.DeleteEntity(&show); err != nil {
		responses.SendError(w, *err)
		return
	}
	responses.SendNoContent(w)
}
