package booking

import (
	"database/sql"
	"log"
	"net/http"
)

type Config interface {
}

/*
Create a primary http.Handler to be served by the application
*/
func NewServer(config Config, logger *log.Logger, db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, db)
	return mux
}
