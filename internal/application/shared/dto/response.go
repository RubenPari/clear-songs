package dto

import (
	"time"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Error represents an error in the API response
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains metadata about the response
type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
}

// Success creates a successful response
func Success(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// ErrorResponse creates an error response
func ErrorResponse(code, message string) *Response {
	return &Response{
		Success: false,
		Error: &Error{
			Code:    code,
			Message: message,
		},
		Meta: &Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// ValidationError creates a validation error response
func ValidationError(message string) *Response {
	return ErrorResponse("VALIDATION_ERROR", message)
}

// InternalError creates an internal server error response
func InternalError(message string) *Response {
	return ErrorResponse("INTERNAL_ERROR", message)
}

// NotFoundError creates a not found error response
func NotFoundError(resource string) *Response {
	return ErrorResponse("NOT_FOUND", resource+" not found")
}

// UnauthorizedError creates an unauthorized error response
func UnauthorizedError() *Response {
	return ErrorResponse("UNAUTHORIZED", "Authentication required")
}
