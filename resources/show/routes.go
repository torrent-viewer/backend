package show

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/shwoodard/jsonapi"
	"github.com/torrent-viewer/backend/database"
	"github.com/torrent-viewer/backend/responses"
	"github.com/torrent-viewer/backend/router"
)

// RegisterHandlers registers the handlers for the Shows endpoints
func RegisterHandlers(r *router.Router) {
	r.AddRoute(router.Route{
		Path:    "/shows/{id:[0-9]+}",
		Handler: ShowsShow,
		Method:  "GET",
		Name:    "shows.show",
	}).AddRoute(router.Route{
		Path:    "/shows/{id:[0-9]+}",
		Handler: ShowsUpdate,
		Method:  "PUT",
		Name:    "shows.update",
	}).AddRoute(router.Route{
		Path:    "/shows/{id:[0-9]+}",
		Handler: ShowsDestroy,
		Method:  "DELETE",
		Name:    "shows.delete",
	}).AddRoute(router.Route{
		Path:    "/shows",
		Handler: ShowsIndex,
		Method:  "GET",
		Name:    "shows.index",
	}).AddRoute(router.Route{
		Path:    "/shows",
		Handler: ShowsStore,
		Method:  "POST",
		Name:    "shows.store",
	})
}

// ShowsIndex is the HTTP endpoint used to list Shows instances
func ShowsIndex(w http.ResponseWriter, r *http.Request) {
	var shows Shows
	if err := database.Conn.Find(&shows).Error; err != nil {
		e := responses.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusInternalServerError, e); err != nil {
			log.Fatal(err)
		}
	}
	serialized := make([]interface{}, len(shows))
	for i, s := range shows {
		serialized[i] = s
	}
	if err := responses.SendEntities(w, serialized); err != nil {
		log.Fatal(err)
	}
}

// ShowsStore is the HTTP endpoint used to create new Shows instances
func ShowsStore(w http.ResponseWriter, r *http.Request) {
	var show Show
	if err := jsonapi.UnmarshalPayload(r.Body, &show); err != nil {
		e := responses.Error{
			ID:     "malformated-input",
			Status: "400",
			Title:  "Malformated input",
			Detail: err.Error(),
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if result, err := govalidator.ValidateStruct(show); err != nil || result != true {
		e := responses.Error{
			ID:     "validation-error",
			Status: "400",
			Title:  "Validation Error",
		}
		if err != nil {
			e.Detail = err.Error()
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if database.Conn.NewRecord(&show) != true {
		e := responses.Error{
			ID:     "duplicate-entry",
			Status: "409",
			Title:  "Duplicate Entry",
			Detail: "Trying to create a resource with an existing ID",
			Source: responses.ErrorSource{
				Pointer: "/data/id",
			},
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusConflict, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := database.Conn.Create(&show).Error; err != nil {
		e := responses.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusInternalServerError, e); err != nil {
			log.Fatal(err)
		}
	}
	headers := map[string]string{
		"Location": fmt.Sprintf("/shows/%d", show.ID),
	}
	if err := responses.SendEntity(w, &show, http.StatusCreated, headers); err != nil {
		log.Fatal(err)
	}
}

// ShowsShow is the HTTP endpoint used to show Shows instance by ID
func ShowsShow(w http.ResponseWriter, r *http.Request) {
	vars := router.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e := responses.Error{
			ID:     "integer-conversion",
			Status: "400",
			Title:  "Integer Conversion Error",
			Detail: err.Error(),
		}
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	var show Show
	d := database.Conn.First(&show, id)
	if d.RecordNotFound() != false {
		err := d.Error
		e := responses.Error{
			ID:     "not-found",
			Status: "404",
			Title:  "Not Found",
			Detail: err.Error(),
		}
		if err := responses.SendError(w, http.StatusNotFound, e); err != nil {
			log.Fatal(err)
		}
		return
	} else if err := d.Error; err != nil {
		e := responses.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
		if err := responses.SendError(w, http.StatusInternalServerError, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := responses.SendEntity(w, &show, http.StatusOK, nil); err != nil {
		log.Fatal(err)
	}
}

// ShowsUpdate is the HTTP endpoint used to update a Show instance by its ID
func ShowsUpdate(w http.ResponseWriter, r *http.Request) {
	vars := router.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e := responses.Error{
			ID:     "integer-conversion",
			Status: "400",
			Title:  "Integer Conversion Error",
			Detail: err.Error(),
		}
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	var show Show
	if err := jsonapi.UnmarshalPayload(r.Body, &show); err != nil {
		e := responses.Error{
			ID:     "malformated-input",
			Status: "400",
			Title:  "Malformated input",
		}
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if show.ID != id {
		e := responses.Error{
			ID:     "unmatched-ids",
			Status: "400",
			Title:  "IDs do not match",
			Detail: "The URL ID does not match the input ID",
			Source: responses.ErrorSource{
				Pointer: "/data/id",
			},
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if result, err := govalidator.ValidateStruct(show); err != nil || result != true {
		e := responses.Error{
			ID:     "validation-error",
			Status: "400",
			Title:  "Validation Error",
		}
		if err != nil {
			e.Detail = err.Error()
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if database.Conn.NewRecord(&show) == true {
		e := responses.Error{
			ID:     "not-found",
			Status: "404",
			Title:  "Not Found",
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusNotFound, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := database.Conn.Model(&show).Update(&show).Error; err != nil {
		e := responses.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusInternalServerError, e); err != nil {
			log.Fatal(err)
		}
	}
	if err := responses.SendNoContent(w); err != nil {
		log.Fatal(err)
	}
}

// ShowsDestroy is the HTTP endpoint used to delete a Show instance by its ID
func ShowsDestroy(w http.ResponseWriter, r *http.Request) {
	vars := router.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		e := responses.Error{
			ID:     "integer-conversion",
			Status: "400",
			Title:  "Integer Conversion Error",
			Detail: err.Error(),
		}
		if err := responses.SendError(w, http.StatusBadRequest, e); err != nil {
			log.Fatal(err)
		}
		return
	}
	show := Show{
		ID: id,
	}
	if err := database.Conn.Delete(&show).Error; err != nil {
		e := responses.Error{
			ID:     "database-error",
			Status: "500",
			Title:  "Database Error",
			Detail: err.Error(),
		}
		log.Println(e)
		if err := responses.SendError(w, http.StatusInternalServerError, e); err != nil {
			log.Fatal(err)
		}
	}
	if err := responses.SendNoContent(w); err != nil {
		log.Fatal(err)
	}
}
