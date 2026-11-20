package failure

import "tegami-server/internal/errs"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Response struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func NewResponse(code, message string, details []FieldError) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewFieldError(appFE errs.FieldError) FieldError {
	return FieldError{
		Field:   appFE.Field,
		Message: appFE.Message,
	}
}
