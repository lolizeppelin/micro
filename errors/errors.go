// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(id, detail string, code int32) error {
	return NewError(id, detail, code)
}

func NewError(id, detail string, code int32) *Error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

// Equal tries to compare errors.
func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

// As finds the first error in err's chain that matches *Error.
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var merr *Error
	if errors.As(err, &merr) {
		return merr, true
	}
	return nil, false
}

func NewMultiError() *MultiError {
	return &MultiError{
		Errors: make([]*Error, 0),
	}
}

func (e *MultiError) Append(err ...*Error) {
	e.Errors = append(e.Errors, err...)
}

func (e *MultiError) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *MultiError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
