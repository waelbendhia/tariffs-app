package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3" //SQLite3 driver
	"github.com/pkg/errors"
)

// OpenSQLite will open SQLite3 database at path, if it does not exist it will
// be created
func OpenSQLite(path string) *sql.DB {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(errors.Wrapf(err, "Could not create file: %s", path))
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(errors.Wrapf(err, "Could not open database at: %s", path))
	}
	return db
}
