package sqlite

import (
	"bitespeed_task/models"
	"database/sql"
	"time"

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

func NewSqliteModel(dsn string) (*ContactsSqlite, error) {
	db, err := getDBConnection(dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(create_contacts_table)
	if err != nil {
		return nil, err
	}

	return &ContactsSqlite{dsn: dsn}, nil
}

func reconciliateWithEmail(tx *sql.Tx, email string) (*sql.Rows, error) {
	var (
		id, linkedId int64
	)
	row := tx.QueryRow(fetch_email_id_and_email_primary_id, email)
	if err := row.Scan(&id, &linkedId); err != nil {
		return nil, err
	}
	// no match
	if id == 0 {
		err := tx.QueryRow(insert_email_when_no_match, email).Scan(&id)
		if err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(fetch_reconciliated_records_given_id, id)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func reconciliateWithPhone(tx *sql.Tx, phone string) (*sql.Rows, error) {
	var (
		id, linkedId int64
	)
	row := tx.QueryRow(fetch_phone_id_and_phone_primary_id, phone)
	if err := row.Scan(&id, &linkedId); err != nil {
		return nil, err
	}
	// no match
	if id == 0 {
		err := tx.QueryRow(insert_phone_when_no_match, phone).Scan(&id)
		if err != nil {
			return nil, err
		}
	}

	rows, err := tx.Query(fetch_reconciliated_records_given_id, id)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func reconciliateWithEmailAndPhone(tx *sql.Tx, email, phone string) (*sql.Rows, error) {

	var (
		email_id, email_primary_id,
		phone_id, phone_primary_id,
		chosen_primary_id int64
		rows *sql.Rows
	)

	err := tx.QueryRow(fetch_email_id_and_email_primary_id, email).Scan(
		&email_id, &email_primary_id,
	)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(fetch_phone_id_and_phone_primary_id, phone).Scan(
		&phone_id, &phone_primary_id,
	)
	if err != nil {
		return nil, err
	}

	switch {
	// no match
	case email_primary_id == 0 && phone_primary_id == 0:
		_, err = tx.Exec(insert_when_email_primary_and_phone_primary_not_exist, email, phone)
		if err != nil {
			return nil, err
		}
		rows, err = tx.Query(fetch_reconciliated_records_given_email_and_phone, email, phone)
		if err != nil {
			return nil, err
		}
	case email_primary_id != 0 && phone_primary_id != 0:
		err = tx.QueryRow(fetch_chosen_primary_id, email_primary_id, phone_primary_id).Scan(&chosen_primary_id)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(update_demote_email_primary_or_phone_primary_to_secondary, chosen_primary_id, email_primary_id, phone_primary_id, chosen_primary_id)
		if err != nil {
			return nil, err
		}
		_, err = tx.Exec(update_linkedId_for_all_followers_of_not_the_chosen_one, chosen_primary_id, email_primary_id, phone_primary_id)
		if err != nil {
			return nil, err
		}
		rows, err = tx.Query(fetch_reconciliated_records_given_id, chosen_primary_id)
		if err != nil {
			return nil, err
		}
	case email_primary_id != 0 && phone_primary_id == 0:
		_, err = tx.Exec(insert_when_only_one_primary_exists, email, phone, email_primary_id)
		if err != nil {
			return nil, err
		}
		rows, err = tx.Query(fetch_reconciliated_records_given_id, email_primary_id, email_primary_id)
		if err != nil {
			return nil, err
		}
	case email_primary_id == 0 && phone_primary_id != 0:
		_, err = tx.Exec(insert_when_only_one_primary_exists, email, phone, phone_primary_id)
		if err != nil {
			return nil, err
		}
		rows, err = tx.Query(fetch_reconciliated_records_given_id, phone_primary_id, phone_primary_id)
		if err != nil {
			return nil, err
		}
	}

	return rows, nil
}

func (m *ContactsSqlite) Reconciliate(email, phone string) ([]*models.Contact, error) {
	db, err := getDBConnection(m.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	var (
		rows *sql.Rows
	)
	switch {
	case email != "" && phone != "":
		rows, err = reconciliateWithEmailAndPhone(tx, email, phone)
	case email != "":
		rows, err = reconciliateWithEmail(tx, email)
	case phone != "":
		rows, err = reconciliateWithPhone(tx, phone)
	}
	if err != nil {
		return nil, err
	}

	defer func() {
		if rows != nil {
			rows.Close()
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		db.Close()
	}()

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

func (m *ContactsSqlite) Close() {
}
