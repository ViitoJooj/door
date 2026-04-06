package database

import "database/sql"

var DB *sql.DB

func Conn() error {
	var err error

	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	return nil
}
