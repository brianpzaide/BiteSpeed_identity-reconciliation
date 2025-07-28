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
	existing_primary_id_query          = `SELECT id FROM contacts WHERE (phoneNumber = ? OR email = ?) AND linkPrecedence = 'primary' ORDER BY created_at LIMIT 1;`
	existing_phone_number_match_query  = `SELECT id, created_at FROM contacts WHERE phoneNumber = ? AND linkPrecedence = 'primary' LIMIT 1;`
	existing_email_match_query         = `SELECT id, created_at FROM contacts WHERE email = ? AND linkPrecedence = 'primary' LIMIT 1;`
	insert_primary_contact_query       = `INSERT INTO contacts (phoneNumber, email, linkPrecedence, createdAt, updatedAt) VALUES (?, ?, 'primary', current_timestamp, current_timestamp);`
	insert_primary_contact_email_query = `INSERT INTO contacts (email, linkPrecedence, createdAt, updatedAt) VALUES (?, 'primary', current_timestamp, current_timestamp);`
	insert_primary_contact_phone_query = `INSERT INTO contacts (phoneNumber, linkPrecedence, createdAt, updatedAt) VALUES (?, 'primary', current_timestamp, current_timestamp);`
	insert_secondary_contact_query     = `INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt) VALUES (?, ?, ?, 'secondary', current_timestamp, current_timestamp);`
)
