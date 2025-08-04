package postgres

import (
	"bitespeed_task/models"
	"context"
	"database/sql"
	"fmt"
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

func NewPostgresModel(dsn string) (*ContactsPostgres, error) {
	db, err := openPostgresDB(dsn)
	if err != nil {
		return nil, err
	}

	// create the contacts table
	_, err = db.Exec(create_contacts_table)
	if err != nil {
		return nil, err
	}

	// create the stored procedure for reconciliating the contacts
	_, err = db.Exec(create_stored_procedure_with_advisory_lock)
	if err != nil {
		return nil, err
	}

	// insert test data
	// _, err = db.Exec(ADD_TEST_DATA)
	// if err != nil {
	// 	fmt.Println("error inserting test data")
	// 	return nil, err
	// }

	return &ContactsPostgres{db: db}, nil
}

func (m *ContactsPostgres) Reconciliate(email, phoneNumber string) ([]*models.Contact, error) {

	var (
		rows *sql.Rows
		err  error
	)

	if email == "" || phoneNumber == "" {
		if email == "" {
			rows, err = m.db.Query(`SELECT * FROM reconcile_contact(NULL, $1);`, phoneNumber)
		} else {
			rows, err = m.db.Query(`SELECT * FROM reconcile_contact($1, NULL);`, email)
		}
	} else {
		rows, err = m.db.Query(`SELECT * FROM reconcile_contact($1, $2);`, email, phoneNumber)
	}

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	contacts := make([]*models.Contact, 0)
	var (
		linkedId  *int64
		deletedAt *time.Time
	)

	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
			&contact.PhoneNumber,
			&contact.Email,
			&linkedId,
			&contact.LinkPrecedence,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}
		if linkedId != nil {
			contact.LinkedId = *linkedId
		}
		if deletedAt != nil {
			contact.DeletedAt = *deletedAt
		}
		contacts = append(contacts, &contact)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (m *ContactsPostgres) Close() {
	m.db.Close()
}
