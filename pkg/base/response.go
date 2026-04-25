package base

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrCode represents a standardized error code for API responses.
type ErrCode int

const (
	// ErrCodeUnknown represents an unknown error.
	ErrCodeUnknown ErrCode = 1000
	// ErrCodeBadRequest represents a bad request error.
	ErrCodeBadRequest ErrCode = 1001
	// ErrCodeNotFound represents a not found error.
	ErrCodeNotFound ErrCode = 1002
	// ErrCodeConflict represents a conflict error.
	ErrCodeConflict ErrCode = 1003
	// ErrCodeInternal represents an internal server error.
	ErrCodeInternal ErrCode = 1004
)

// GopeedError represents a structured error returned by the Gopeed API.
type GopeedError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e *GopeedError) Error() string {
	return fmt.Sprintf("gopeed error (code=%d): %s", e.Code, e.Message)
}

// NewGopeedError creates a new GopeedError with the given code and message.
func NewGopeedError(code ErrCode, msg string) *GopeedError {
	return &GopeedError{
		Code:    code,
		Message: msg,
	}
}

// Response is the generic API response wrapper used across Gopeed HTTP handlers.
type Response[T any] struct {
	Code ErrCode `json:"code"`
	Msg  string  `json:"msg,omitempty"`
	Data T       `json:"data,omitempty"`
}

// Ok returns a successful Response containing the provided data.
func Ok[T any](data T) *Response[T] {
	return &Response[T]{
		Code: 0,
		Data: data,
	}
}

// Err returns an error Response with the given error code and message.
func Err(code ErrCode, msg string) *Response[any] {
	return &Response[any]{
		Code: code,
		Msg:  msg,
	}
}

// WriteJSON writes a JSON-encoded response to the http.ResponseWriter with the
// appropriate Content-Type header and the given HTTP status code.
func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

// WriteOk writes a successful JSON response (HTTP 200) wrapping the given data.
func WriteOk[T any](w http.ResponseWriter, data T) error {
	return WriteJSON(w, http.StatusOK, Ok(data))
}

// WriteErr writes an error JSON response with the given HTTP status code,
// error code, and message.
func WriteErr(w http.ResponseWriter, statusCode int, code ErrCode, msg string) error {
	return WriteJSON(w, statusCode, Err(code, msg))
}

// WriteGopeedError writes a GopeedError as a JSON response, mapping the error
// code to an appropriate HTTP status code.
func WriteGopeedError(w http.ResponseWriter, err *GopeedError) error {
	httpStatus := gopeedErrToHTTPStatus(err.Code)
	return WriteErr(w, httpStatus, err.Code, err.Message)
}

// gopeedErrToHTTPStatus maps a GopeedError code to a standard HTTP status code.
func gopeedErrToHTTPStatus(code ErrCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeInternal:
		return http.StatusInternalServerError
	default:
		// ErrCodeUnknown and any unrecognized codes fall back to 500.
		return http.StatusInternalServerError
	}
}
