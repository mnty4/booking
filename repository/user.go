package repository

import (
	"database/sql"

	"github.com/mnty4/booking/model"
)

func InsertUser(db *sql.DB, user model.User) (id int64, err error) {
	query := "INSERT INTO users (email, first_name, last_name) VALUES (?, ?, ?)"
	res, err := db.Exec(query, user.Email, user.FirstName, user.LastName)
	if err != nil {
		return 0, err
	}
	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
