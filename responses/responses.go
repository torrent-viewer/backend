package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shwoodard/jsonapi"
)

type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

// Error represent an API Error that will be sent to a client
type Error struct {
	ID    string `json:"id"`
	Links struct {
		About string `json:"about,omitempty"`
	} `json:"links,omitempty"`
	Status string      `json:"status,omitempty"`
	Code   string      `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Source ErrorSource `json:"source,omitempty"`
	Meta   interface{} `json:"meta,omitempty"`
}

type Errors []Error

func (e Error) Error() string {
	return fmt.Sprintf("HTTP %s: %s (%s)", e.Code, e.Title, e.ID)
}

type errorResponse struct {
	Errors Errors `json:"errors"`
}

// SendError writes a single Error to w
func SendError(w http.ResponseWriter, errcode int, e Error) error {
	w.WriteHeader(errcode)
	response := errorResponse{
		Errors: Errors{
			e,
		},
	}
	return json.NewEncoder(w).Encode(response)
}

// SendErrors writes multiple Errors to w
func SendErrors(w http.ResponseWriter, errcode int, e Errors) error {
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
