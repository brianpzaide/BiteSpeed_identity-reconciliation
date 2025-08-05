package sqlite

const create_contacts_table = `CREATE TABLE IF NOT EXISTS contacts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phoneNumber TEXT,
    email TEXT,
    linkedId INTEGER REFERENCES contacts(id),
    linkPrecedence TEXT CHECK (linkPrecedence IN ('primary', 'secondary')),
    createdAt INTEGER NOT NULL,
    updatedAt INTEGER NOT NULL,
    deletedAt INTEGER
)`

const (
	// when only email is provided
	insert_email_when_no_match = `INSERT INTO contacts (email, linkPrecedence) VALUES (?, 'primary') RETURNING id;`

	// when only phone is provided
	insert_phone_when_no_match = `INSERT INTO contacts (phoneNumber, linkPrecedence) VALUES (?, 'primary') RETURNING id;`

	// When both email and phone are provided
	fetch_reconciliated_records_given_email_and_phone         = `SELECT * FROM contacts where email = ? AND phoneNumber = ? ORDER By createdAt;`
	fetch_chosen_primary_id                                   = `SELECT id FROM contacts WHERE id IN (?, ?) ORDER BY createdAt ASC LIMIT 1;`
	insert_when_email_primary_and_phone_primary_not_exist     = `INSERT INTO contacts (email, phoneNumber, linkPrecedence) VALUES (?, ?, 'primary')`
	fetch_when_email_primary_and_phone_primary_not_exist      = `SELECT * FROM CONTACTS WHERE email = ? AND phoneNumber = ? ORDER BY createdAt LIMIT 1;`
	update_demote_email_primary_or_phone_primary_to_secondary = `UPDATE contacts SET linkedId = ?, linkPrecedence = 'secondary', updatedAt = current_timestamp WHERE id IN (?, ?) AND id != ?;`
	update_linkedId_for_all_followers_of_not_the_chosen_one   = `UPDATE contacts SET linkedId = ?, updatedAt = current_timestamp WHERE linkedId IN (?, ?);`
	insert_when_only_one_primary_exists                       = `INSERT INTO contacts (email, phoneNumber, linkedId, linkPrecedence) VALUES (?, ?, ?, 'secondary');`

	fetch_email_id_and_email_primary_id  = `SELECT id, COALESCE(linkedId, id) FROM contacts WHERE email = ? ORDER BY createdAt LIMIT 1;`
	fetch_phone_id_and_phone_primary_id  = `SELECT id, COALESCE(linkedId, id) FROM contacts WHERE phoneNumber = ? ORDER BY createdAt LIMIT 1;`
	fetch_reconciliated_records_given_id = `SELECT * FROM contacts where id = ? OR linkedId = ? ORDER By createdAt;`
)
