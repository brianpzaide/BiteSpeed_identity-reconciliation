package sqlite

import (
	"bitespeed_task/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ContactsSqlite struct {
	dsn string
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

func NewSqliteModel(dsn string) (ContactsSqlite, error) {
	db, err := getDBConnection(dsn)
	if err != nil {
		return ContactsSqlite{}, err
	}
	defer db.Close()

	_, err = db.Exec(create_contacts_table)
	if err != nil {
		return ContactsSqlite{}, err
	}

	return ContactsSqlite{dsn: dsn}, nil
}

func (m ContactsSqlite) Reconciliate(email, phone string) ([]*models.Contact, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return nil, nil
}

func (m ContactsSqlite) Close() {
}
