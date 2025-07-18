package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ContactsSqlite struct {
}

func getDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewSqliteModel() (ContactsSqlite, error) {
	return ContactsSqlite{}, nil
}

func (c *ContactsSqlite) Fetch(email, phone string) {

}
