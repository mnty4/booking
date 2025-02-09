package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mnty4/booking/api"
)

type Config interface {
}

/*
Create a primary http.Handler to be served by the application
*/
func NewServer(getEnv func(string) string, config Config, logger *log.Logger, db *sql.DB) http.Server {
	mux := http.NewServeMux()
	api.AddRoutes(mux, db, logger)
	return http.Server{
		Handler: mux,
		Addr:    ":" + getEnv("PORT"),
	}
}
