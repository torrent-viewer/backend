package herr

import (
	"fmt"
	"strconv"
)

// ErrorSource point to the source of an error.
// A pointer to the invalid data or the name of the parameter that caused the error.
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

var UnmatchingIDsError = Error{
	ID:     "unmatching-ids",
	Status: "400",
	Title:  "IDs do not match",
	Detail: "The URL ID does not match the input ID",
	Source: ErrorSource{
		Pointer: "/data/id",
	},
}

var DuplicateEntryError = Error{
	ID:     "duplicate-entry",
	Status: "409",
	Title:  "Duplicate Entry",
	Detail: "Trying to create a resource with an existing ID",
	Source: ErrorSource{
		Pointer: "/data/id",
	},
}

func (e Error) Error() string {
	return fmt.Sprintf("HTTP %s: %s (%s)", e.Code, e.Title, e.ID)
}

func (e Error) StatusCode() int {
	status, err := strconv.Atoi(e.Status)
	if err != nil {
		return 500
	}
	return status
}