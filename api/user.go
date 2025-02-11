package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mnty4/booking/errutil"
	"github.com/mnty4/booking/model"
)

func UserPostHandler(db *sql.DB, logger *log.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Printf("Error parsing request body from JSON: %v\n", err)
			if err := errutil.WriteAPIError(w, "Error parsing JSON request body.", 400, "BAD_REQUEST", nil); err != nil {
				logger.Printf("Error writing APIError: %v\n", err)
			}
			return
		}
		r.Body.Close()
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			if err := errutil.WriteValidationError(w, err); err != nil {
				logger.Printf("Error writing ValidationError: %v\n", err)
			}
			return
		}

		// w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// w.Header().Set("Location", fmt.Sprintf("%s/users/%s", id))
	}
}
