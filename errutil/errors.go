package errutil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// The base error struct for sending errors to a client
// - Message: A human-readable description of the error.
// - Code: A high-level error code, defaulted to the HTTP status code.
// - Status: A string for programmatic error classification; more granular than the code.
// - Details: Additional custom error information unique to the endpoint.
type APIError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	status  string
	Details []interface{} `json:"details"`
}

func NewAPIError(message string, code int, status string, details []interface{}) APIError {
	apiError := APIError{}
	apiError.SetStatus(status)
	apiError.Code = code
	apiError.Message = message
	apiError.Details = details
	return apiError
}
func (e *APIError) GetStatus() string {
	return e.status
}
func (e *APIError) SetStatus(status string) {
	switch status {
	case "INTERNAL":
		break
	case "BAD_REQUEST":
		break
	case "VALIDATION":
		break
	default:
		status = "UNKNOWN"
	}
	e.status = status
}

func WriteAPIError(w http.ResponseWriter, message string, code int, status string, details []interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	apiError := NewAPIError(message, code, status, details)
	if err := json.NewEncoder(w).Encode(&apiError); err != nil {
		return err
	}
	return nil
}

func WriteBadRequestError(w http.ResponseWriter, message string) error {
	return WriteAPIError(w, message, http.StatusBadRequest, "BAD_REQUEST", nil)
}

func WriteInternalError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	apiError := NewAPIError(
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
		"INTERNAL",
		nil,
	)
	if err := json.NewEncoder(w).Encode(&apiError); err != nil {
		return err
	}
	return nil
}

func WriteValidationError(w http.ResponseWriter, err error) error {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return WriteInternalError(w)
	}
	var message string
	details := make([]string, 0)
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range err {
			details = append(details, fmt.Sprintf("%s: %s", validationErr.Field(), validationErr.Error()))
		}
	} else {
		message = "some message"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	apiError := NewAPIError(message, http.StatusBadRequest, "VALIDATION", []interface{}{details})
	if err := json.NewEncoder(w).Encode(&apiError); err != nil {
		return err
	}
	return nil
}
