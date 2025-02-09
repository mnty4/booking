package errutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewAPIError(w http.ResponseWriter, logger *log.Logger, code int, err error, message string) {
	logger.Printf("Error %d: %v\n", code, err)
	apiError := APIError{
		Code:    code,
		Message: message,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(apiError.Code)
	if err := json.NewEncoder(w).Encode(apiError); err != nil {
		logger.Printf("Error encoding APIError: %v", err)
	}
}

func NewInternalError(w http.ResponseWriter, logger *log.Logger, code int, err error) {
	NewAPIError(w, logger, http.StatusInternalServerError, err, "Internal Server error.")
}

func WriteValidationError(w http.ResponseWriter, logger *log.Logger, err error) {
	if err, ok := err.(*validator.InvalidValidationError); ok {
		logger.Println(*err)
		WriteInternalError(w, logger, err)
	}
	w.WriteHeader(400)
	if err, ok := err.(validator.ValidationErrors); ok {
		logger.Println(err)
		for _, validationErr := range err {
			w.Write([]byte(fmt.Sprintf("%s: %s", validationErr.Field(), validationErr.Error())))
		}
	}
}
