package requests

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/shwoodard/jsonapi"
	"github.com/torrent-viewer/backend/router"
	"github.com/torrent-viewer/backend/herr"
)

func ParseID(r *http.Request) (int, *herr.Error) {
	vars := router.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, &herr.Error{
			ID:     "integer-conversion",
			Status: "400",
			Title:  "Integer Conversion Error",
			Detail: err.Error(),
			Source: herr.ErrorSource{
				Parameter: "id",
			},
		}
	}
	return id, nil
}

func ReceiveEntity(r *http.Request, entity interface{}) *herr.Error {
	if err := jsonapi.UnmarshalPayload(r.Body, entity); err != nil {
		return &herr.Error{
			ID:     "malformated-input",
			Status: "400",
			Title:  "Malformated input",
			Detail: err.Error(),
		}
	}
	if result, err := govalidator.ValidateStruct(entity); err != nil || result != true {
		e := herr.Error{
			ID:     "validation-error",
			Status: "400",
			Title:  "Validation Error",
		}
		if err != nil {
			e.Detail = err.Error()
		}
		return &e
	}
	return nil
}