package db

import (
	"database/sql"
	"fmt"
)

func Connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		return db, fmt.Errorf("Error opening DB: %v", err)
	}

	return db, nil
}
