package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func AddRoutes(logger *log.Logger, db *sql.DB, validate *validator.Validate, mux *http.ServeMux) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("POST /api/users", UserCreateHandler(logger, db, validate))
	// mux.HandleFunc("GET /api/bookings", bookingGetHandler())
	// mux.HandleFunc("POST /api/bookings", bookingPostHandler())
}
