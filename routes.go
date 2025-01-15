package booking

import (
	"database/sql"
	"net/http"
)

func addRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.Handle("/", http.NotFoundHandler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	// mux.HandleFunc("GET /api/bookings", bookingGetHandler())
	// mux.HandleFunc("POST /api/bookings", bookingPostHandler())
}
