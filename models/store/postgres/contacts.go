package postgres

import (
	"bitespeed_task/models"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type ContactsPostgres struct {
	db *sql.DB
}

func openPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresModel(dsn string) (ContactsPostgres, error) {
	db, err := openPostgresDB(dsn)
	if err != nil {
		return ContactsPostgres{}, err
	}

	// create the contacts table
	_, err = db.Exec(create_contacts_table)
	if err != nil {
		return ContactsPostgres{}, err
	}

	// create the stored procedure for reconciliating the contacts
	_, err = db.Exec(create_stored_procedure_with_advisory_lock)
	if err != nil {
		return ContactsPostgres{}, err
	}

	return ContactsPostgres{db: db}, nil
}

func (m ContactsPostgres) Reconciliate(email, phoneNumber string) ([]*models.Contact, error) {
	return nil, nil
}

func (m ContactsPostgres) Close() {
	m.db.Close()
}
