package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mnty4/booking/errutil"
	"github.com/mnty4/booking/model"
	"github.com/mnty4/booking/repository"
)

func UserCreateHandler(logger *log.Logger, db *sql.DB, validate *validator.Validate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			logger.Printf("Error parsing request body from JSON: %v\n", err)
			if err := errutil.WriteBadRequestError(w, "Error parsing JSON request body."); err != nil {
				logger.Printf("Error writing BadRequestError: %v\n", err)
			}
			return
		}
		r.Body.Close()
		if err := validate.Struct(user); err != nil {
			logger.Printf("Error validating user: %v\n", err)
			if err := errutil.WriteValidationError(w, err); err != nil {
				logger.Printf("Error writing ValidationError: %v\n", err)
			}
			return
		}
		id, err := repository.InsertUser(db, user)
		if err != nil {
			logger.Printf("Error inserting user: %v\n", err)
			if err := errutil.WriteBadRequestError(w, "Error creating user."); err != nil {
				logger.Printf("Error writing BadRequestError: %v\n", err)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Location", fmt.Sprintf("%s/api/users/%d", os.Getenv("BASE_URL"), id))
	}
}
