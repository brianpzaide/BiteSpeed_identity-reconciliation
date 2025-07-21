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
