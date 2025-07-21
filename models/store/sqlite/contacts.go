package sqlite

import (
	"bitespeed_task/models"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	existing_primary_id_query          = `SELECT id FROM contacts WHERE (phoneNumber = ? OR email = ?) AND linkPrecedence = 'primary' ORDER BY created_at LIMIT 1;`
	existing_phone_number_match_query  = `SELECT id, created_at FROM contacts WHERE phoneNumber = ? AND linkPrecedence = 'primary' LIMIT 1;`
	existing_email_match_query         = `SELECT id, created_at FROM contacts WHERE email = ? AND linkPrecedence = 'primary' LIMIT 1;`
	insert_primary_contact_query       = `INSERT INTO contacts (phoneNumber, email, linkPrecedence, createdAt, updatedAt) VALUES (?, ?, 'primary', current_timestamp, current_timestamp);`
	insert_primary_contact_email_query = `INSERT INTO contacts (email, linkPrecedence, createdAt, updatedAt) VALUES (?, 'primary', current_timestamp, current_timestamp);`
	insert_primary_contact_phone_query = `INSERT INTO contacts (phoneNumber, linkPrecedence, createdAt, updatedAt) VALUES (?, 'primary', current_timestamp, current_timestamp);`
	insert_secondary_contact_query     = `INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES (?, ?, ?, 'secondary', current_timestamp, current_timestamp);`
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

func fetch_existing_primary_id(tx *sql.Tx, email, phone string) (int64, error) {
	query := existing_primary_id_query
	var id int64
	if email == "" {
		query = `SELECT id FROM contacts WHERE phoneNumber = ? AND linkPrecedence = 'primary' ORDER BY created_at LIMIT 1`
		err := tx.QueryRow(query, phone).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil
			} else {
				return 0, err
			}
		} else {
			return id, nil
		}
	}
	if phone == "" {
		query = `SELECT id FROM contacts WHERE email = ? AND linkPrecedence = 'primary' ORDER BY created_at LIMIT 1`
		err := tx.QueryRow(query, email).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, nil
			} else {
				return 0, err
			}
		} else {
			return id, nil
		}
	}

	err := tx.QueryRow(query, phone, email).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		} else {
			return 0, err
		}
	} else {
		return id, nil
	}
}

func fetch_existing_email_match(tx *sql.Tx, email string) (*sql.Row, error) {
	return nil, nil
}

func fetch_existing_phone_match(tx *sql.Tx, phone string) (*sql.Row, error) {
	return nil, nil
}

func insert_primary_contact(tx *sql.Tx, email, phone string) error {
	if email == "" {
		_, err := tx.Exec(insert_primary_contact_query, phone)
		return err
	}
	if phone == "" {
		_, err := tx.Exec(insert_primary_contact_query, email)
		return err
	}
	_, err := tx.Exec(insert_primary_contact_query, phone, email)
	return err
}

func (m ContactsSqlite) reconciliate(email, phone string) ([]*models.Contact, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	existingPrimary_id, err := fetch_existing_primary_id(tx, email, phone)
	if err != nil {
		return nil, err
	}
	if existingPrimary_id == 0 {
		err = insert_primary_contact(tx, email, phone)
		if err != nil {
			return nil, err
		}
	} else {

	}

	// existingPrimary_id, err := fetch_existing_primary_id(tx, email, phone)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (m ContactsSqlite) Reconciliate(email, phone string) ([]*models.Contact, error) {

	return nil, nil
}

func (m ContactsSqlite) Close() {
}
