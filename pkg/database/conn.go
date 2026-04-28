package database

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func hasProtocolApplyScopeColumn(db *sql.DB) bool {
	rows, err := db.Query("PRAGMA table_info(protocol_settings);")
	if err != nil {
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var colType string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &colType, &notNull, &defaultValue, &pk); err != nil {
			return false
		}
		if name == "apply_scope" {
			return true
		}
	}
	return false
}

func Conn() {
	var err error

	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	driver, err := sqlite3.WithInstance(DB, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"sqlite3", driver,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if strings.Contains(err.Error(), "Dirty database version 12") {
			version, dirty, versionErr := m.Version()
			if versionErr == nil && dirty && version == 12 {
				if !hasProtocolApplyScopeColumn(DB) {
					panic(err)
				}
				if err := m.Force(12); err != nil {
					panic(err)
				}
				if err := m.Up(); err != nil && err != migrate.ErrNoChange {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
