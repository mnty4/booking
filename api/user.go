package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
			if _, ok := err.(*json.MarshalerError); ok {
				errutil.NewAPIError(w, logger, http.StatusBadRequest, fmt.Errorf("error when decoding json into: %+v: %v", user, err), "Bad JSON")
			} else {
				errutil.NewInternalError(w, logger, http.StatusInternalServerError,
					fmt.Errorf("error when decoding json into: %+v: %v", user, err))
			}
			return
		}
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			errutil.WriteValidationError(w, logger, err)
			return
		}
		r.Body.Close()
		w.WriteHeader(http.StatusCreated)
	}
}
