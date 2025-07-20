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

const reconciliate_in_transaction = `BEGIN IMMEDIATE;

WITH existing_primary AS (
    SELECT id FROM items WHERE some_key = ? AND linked_id IS NULL LIMIT 1
),
inserted AS (
    INSERT INTO items (some_key, linked_id)
    VALUES (?, (SELECT id FROM existing_primary))
    RETURNING *
)
SELECT * FROM inserted;

COMMIT;`
