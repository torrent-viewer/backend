package show

import (
	"fmt"
	"log"
	"net/http"

	"github.com/torrent-viewer/backend/datastore"
	"github.com/torrent-viewer/backend/responses"
	"github.com/torrent-viewer/backend/requests"
	"github.com/torrent-viewer/backend/herr"
)

// ShowsList is the HTTP endpoint used to create list Shows instances
func (res ShowResource) RouteList(w http.ResponseWriter, r *http.Request) {
	var entries Shows
	if err := datastore.FetchEntities(&entries); err != nil {
		if e := responses.SendError(w, http.StatusInternalServerError, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	serialized := make([]interface{}, len(entries), len(entries))
	for i, e := range entries {
		serialized[i] = e
	}
	if err := responses.SendEntities(w, serialized); err != nil {
		log.Fatal(err)
	}
}

// ShowsStore is the HTTP endpoint used to create new Shows instances
func (res ShowResource) RouteStore(w http.ResponseWriter, r *http.Request) {
	var show Show
	if err := requests.ReceiveEntity(r, &show); err != nil {
		if e := responses.SendError(w, http.StatusBadRequest, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	if datastore.Conn.NewRecord(&show) != true {
		if e := responses.SendError(w, http.StatusConflict, herr.DuplicateEntryError); e != nil {
			log.Fatal(e)
		}
		return
	}
	if err := datastore.StoreEntity(&show); err != nil {
		if e := responses.SendError(w, http.StatusInternalServerError, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/shows/%d", show.ID))
	if err := responses.SendEntity(w, &show, http.StatusCreated); err != nil {
		log.Fatal(err)
	}
}

// ShowsView is the HTTP endpoint used to show Shows instance by ID
func (res ShowResource) RouteView(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		if e := responses.SendError(w, http.StatusBadRequest, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	var show Show
	if err := datastore.FetchEntity(&show, id); err != nil {
		if e := responses.SendError(w, http.StatusNotFound, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	if err := responses.SendEntity(w, &show, http.StatusOK); err != nil {
		log.Fatal(err)
	}
}

// ShowsUpdate is the HTTP endpoint used to update a Show instance by its ID
func (res ShowResource) RouteUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		if e := responses.SendError(w, http.StatusBadRequest, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	var show Show
	if err := datastore.FetchEntity(&show, id); err != nil {
		if e := responses.SendError(w, http.StatusNotFound, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	if err := requests.ReceiveEntity(r, &show); err != nil {
		if e := responses.SendError(w, http.StatusBadRequest, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	if show.ID != id {
		if err := responses.SendError(w, http.StatusBadRequest, herr.UnmatchingIDsError); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := datastore.UpdateEntity(&show); err != nil {
		if e := responses.SendError(w, http.StatusInternalServerError, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	if err := responses.SendNoContent(w); err != nil {
		log.Fatal(err)
	}
}

// ShowsDestroy is the HTTP endpoint used to delete a Show instance by its ID
func (res ShowResource) RouteDestroy(w http.ResponseWriter, r *http.Request) {
	id, err := requests.ParseID(r)
	if err != nil {
		if e := responses.SendError(w, http.StatusBadRequest, *err); e != nil {
			log.Fatal(e)
		}
		return
	}
	show := Show{
		ID: id,
	}
	if err := datastore.DeleteEntity(&show); err != nil {
		if e := responses.SendError(w, http.StatusInternalServerError, *err); e != nil {
			log.Fatal(e)
		}
	}
	if err := responses.SendNoContent(w); err != nil {
		log.Fatal(err)
	}
}
