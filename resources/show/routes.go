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

// ShowsList is the HTTP endpoint used to create list Shows instances
func (res ShowResource) RouteList(w http.ResponseWriter, r *http.Request) {
	var entries Shows
	if err := database.Conn.Find(&entries).Error; err != nil {
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
		return
	}
	if err := responses.SendEntities(w, entries); err != nil {
		log.Fatal(err)
	}
}

// ShowsStore is the HTTP endpoint used to create new Shows instances
func (res ShowResource) RouteStore(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Location", fmt.Sprintf("/shows/%d", show.ID))
	if err := responses.SendEntity(w, &show, http.StatusCreated); err != nil {
		log.Fatal(err)
	}
}

// ShowsView is the HTTP endpoint used to show Shows instance by ID
func (res ShowResource) RouteView(w http.ResponseWriter, r *http.Request) {
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
	if err := responses.SendEntity(w, &show, http.StatusOK); err != nil {
		log.Fatal(err)
	}
}

// ShowsUpdate is the HTTP endpoint used to update a Show instance by its ID
func (res ShowResource) RouteUpdate(w http.ResponseWriter, r *http.Request) {
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
func (res ShowResource) RouteDestroy(w http.ResponseWriter, r *http.Request) {
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
