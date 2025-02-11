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
	Message string        `json:"message"`
	Code    int           `json:"code"`
	Status  ErrorStatus   `json:"status"`
	Details []interface{} `json:"details"`
}
type ErrorStatus string

const (
	StatusInternal   ErrorStatus = "INTERNAL"
	StatusBadRequest ErrorStatus = "BAD_REQUEST"
	StatusValidation ErrorStatus = "VALIDATION"
)

func NewAPIError(message string, code int, status ErrorStatus, details []interface{}) APIError {
	apiError := APIError{}
	apiError.Status = status
	apiError.Code = code
	apiError.Message = message
	apiError.Details = details
	return apiError
}

func WriteAPIError(w http.ResponseWriter, message string, code int, status ErrorStatus, details []interface{}) error {
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
