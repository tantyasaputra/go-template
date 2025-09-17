package http

import (
	"encoding/json"
	"go-template/internal/log"
	"net/http"
)

// ErrorField represents error of field
type ErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents the default error response
type ErrorResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Fields  []*ErrorField `json:"fields,omitempty"`
}

// ResponseJSON writes json http response
func ResponseJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorw("Http Error Response", "event", "http response encode error", "service", "IKN_B2B", "error", err)
	}
}

// ResponseError writes error http response
func ResponseError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var errorCode string

	switch status {
	case http.StatusUnauthorized:
		errorCode = "Unauthorized Request"
	case http.StatusNotFound:
		errorCode = "NotFound"
	case http.StatusBadRequest:
		errorCode = "BadRequest"
	case http.StatusUnprocessableEntity:
		errorCode = "UnprocessableEntity"
	case http.StatusTooManyRequests:
		errorCode = "TooManyRequests"
	default:
		errorCode = "InternalServerError"
	}

	log.Errorw("Http Error Response", "event", "http response", "error", err)

	var encodeErr error

	if status == http.StatusInternalServerError {
		encodeErr = json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: "Server error",
		})
	} else {
		encodeErr = json.NewEncoder(w).Encode(ErrorResponse{
			Code:    errorCode,
			Message: err.Error(),
		})
	}

	if encodeErr != nil {
		log.Errorw("Http Error Response", "event", "http response encode error", "service", "IKN_B2B", "error", encodeErr)
	}
}
