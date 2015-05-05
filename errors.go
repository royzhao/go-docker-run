package main

import (
	"encoding/xml"
	"fmt"
)

const (
	// Error codes
	ErrWrongInput = 1
	ErrRunCmd     = 2
	CommitOk      = 3
)

// The serializable Error structure.
type Error struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Status  int      `json:"status" xml:"status,attr"`
	Message string   `json:"message" xml:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}

// NewError creates an error instance with the specified code and message.
func NewError(code int, msg string) *Error {
	return &Error{
		Status:  code,
		Message: msg,
	}
}
