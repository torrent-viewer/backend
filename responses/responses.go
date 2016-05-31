package responses

import (
	"encoding/json"
	"net/http"

	"github.com/shwoodard/jsonapi"
	"github.com/torrent-viewer/backend/herr"
)

type errorResponse struct {
	Errors herr.Errors `json:"errors"`
}

// SendError writes a single Error to w
func SendError(w http.ResponseWriter, e herr.Error) error {
	w.WriteHeader(e.StatusCode())
	response := errorResponse{
		Errors: herr.Errors{
			e,
		},
	}
	return json.NewEncoder(w).Encode(response)
}

// SendErrors writes multiple Errors to w
func SendErrors(w http.ResponseWriter, errcode int, e herr.Errors) error {
	w.WriteHeader(errcode)
	response := errorResponse{
		Errors: e,
	}
	return json.NewEncoder(w).Encode(response)
}

// SendEntity marshalls the given entity and writes it to w
func SendEntity(w http.ResponseWriter, entity interface{}, status int) error {
	w.WriteHeader(status)
	return jsonapi.MarshalOnePayload(w, entity)
}

// SendEntities marshalls the given entities and writes them to w
func SendEntities(w http.ResponseWriter, entities []interface{}) error {
	w.WriteHeader(http.StatusOK)
	return jsonapi.MarshalManyPayload(w, entities)
}

// SendNoContent sends a HTTP 204 to the client
func SendNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
