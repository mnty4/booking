package api

import (
	"database/sql"
	"log"
	"net/http"
)

func AddRoutes(mux *http.ServeMux, db *sql.DB, logger *log.Logger) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mux.HandleFunc("POST /api/users", UserPostHandler(db, logger))
	// mux.HandleFunc("GET /api/bookings", bookingGetHandler())
	// mux.HandleFunc("POST /api/bookings", bookingPostHandler())
}
