package requests

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/shwoodard/jsonapi"
	"github.com/torrent-viewer/backend/datastore"
	"github.com/torrent-viewer/backend/router"
	"github.com/torrent-viewer/backend/herr"
)

var (
	defaultPageSize int = 50
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

type Pagination struct {
	Offset int
	Limit int
}

func Paginate(model interface{}, r *http.Request) (Pagination, *herr.Error) {
	var total int
	if err := datastore.CountEntities(model, &total, nil); err != nil {
		return Pagination{}, err
	}
	queries := r.URL.Query()
	size := defaultPageSize
	offset := 0
	if sizeq, ok := queries["page[size]"]; ok {
		size, err := strconv.Atoi(sizeq[0])
		if err != nil || size <= 0 {
			return Pagination{}, &herr.Error{
				ID:     "invalid-parameter",
				Status: "400",
				Title:  "Invalid query parameter",
				Source: herr.ErrorSource{
					Parameter: "page[size]",
				},
			}
		}		
	}
	if pageq, ok := queries["page[number]"]; ok {
		page, err := strconv.Atoi(pageq[0])
		offset = (page - 1) * size;
		if err != nil || page < 1 || offset > total {
			return Pagination{}, &herr.Error{
				ID:     "invalid-parameter",
				Status: "400",
				Title:  "Invalid query parameter",
				Source: herr.ErrorSource{
					Parameter: "page[number]",
				},
			}
		}
	}
	return Pagination{
		Offset: offset,
		Limit:  size,
	}, nil
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