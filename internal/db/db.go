package db

import (
	"database/sql"
	"log"
)

func NewDb(driverName, dsn string) *sql.DB {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
