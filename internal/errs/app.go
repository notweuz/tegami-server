package errs

import (
	"errors"
	"fmt"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Kind    Kind
	Code    string
	Message string
	Fields  []FieldError
	Err     error
}

func New(kind Kind, code, msg string, err error) *AppError {
	return &AppError{
		Kind:    kind,
		Code:    code,
		Message: msg,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) WithErr(err error) *AppError {
	cp := *e
	cp.Err = err
	return &cp
} 

func (e *AppError) Unwrap() error {
	return e.Err
}

func IsKind(err error, k Kind) bool {
	var ae *AppError
	return errors.As(err, &ae) && ae.Kind == k
}