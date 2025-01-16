package db

import (
	"database/sql"
	"fmt"
	"log"
)

func NewDb(logger *log.Logger, driverName, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging db: %v", err)
	}
	return db, nil
}
