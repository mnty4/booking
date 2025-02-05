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
func NewServer(getEnv func(string) string, config Config, logger *log.Logger, db *sql.DB) http.Server {
	mux := http.NewServeMux()
	addRoutes(mux, db)
	return http.Server{
		Handler: mux,
		Addr:    ":" + getEnv("PORT"),
	}
}
