package middleware

import (
	"net/http"

	"github.com/goccy/go-json"
)

// ErrorResponse is a custom error response struct for apis
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

// Error makes it compatible with the `error` interface.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new Error instance with an optional message
func NewError(code int, message ...string) *Error {
	err := &Error{
		Code:    code,
		Message: http.StatusText(code),
	}
	if len(message) > 0 {
		err.Message = message[0]
	}
	return err
}

// Error represents an error that occurred while handling a request.
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Map is a shortcut for map[string]interface{}, useful for JSON returns
type Map map[string]interface{}

//JSON response formatting
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
	w.WriteHeader(status)
	w.Write(jsonResponse)
}
