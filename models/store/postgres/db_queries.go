package postgres

const create_contacts_table = `CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    phoneNumber TEXT,
    email TEXT,
    linkedId INTEGER REFERENCES contacts(id),
    linkPrecedence TEXT CHECK (linkPrecedence IN ('primary', 'secondary')),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP
)`

const create_stored_procedure_with_advisory_lock = `CREATE OR REPLACE FUNCTION reconcile_contact(p_email TEXT, p_phone TEXT)
RETURNS TABLE (
    id INTEGER,
    email TEXT,
    phoneNumber TEXT,
    linkedId INTEGER,
    linkPrecedence TEXT,
    createdAt TIMESTAMP,
    updatedAt TIMESTAMP,
	deletedAt TIMESTAMP
)
LANGUAGE plpgsql AS
$$
DECLARE
    email_id INTEGER;
    phone_id INTEGER;
    email_primary_id INTEGER;
    phone_primary_id INTEGER;
    chosen_primary_id INTEGER;
BEGIN
    SELECT contacts.id, COALESCE(contacts.linkedId, contacts.id)
    INTO email_id, email_primary_id
    FROM contacts
    WHERE contacts.email = p_email
    ORDER BY contacts.createdAt
    LIMIT 1;

    SELECT contacts.id, COALESCE(contacts.linkedId, contacts.id)
    INTO phone_id, phone_primary_id
    FROM contacts
    WHERE contacts.phoneNumber = p_phone
    ORDER BY contacts.createdAt
    LIMIT 1;

    IF email_primary_id IS NULL AND phone_primary_id IS NULL THEN
        INSERT INTO contacts (email, phoneNumber, linkPrecedence)
        VALUES (p_email, p_phone, 'primary');

        RETURN QUERY
        SELECT 
            contacts.id, contacts.email, contacts.phoneNumber, contacts.linkedId,
            contacts.linkPrecedence, contacts.createdAt, contacts.updatedAt, contacts.deletedAt
        FROM contacts
        WHERE contacts.email = p_email AND contacts.phoneNumber = p_phone
        ORDER BY contacts.createdAt;
        RETURN;
    END IF;

    IF email_primary_id IS NOT NULL AND phone_primary_id IS NOT NULL THEN
        SELECT contacts.id INTO chosen_primary_id
        FROM contacts
        WHERE contacts.id IN (email_primary_id, phone_primary_id)
        ORDER BY contacts.createdAt ASC
        LIMIT 1;

        UPDATE contacts
        SET linkedId = chosen_primary_id,
            linkPrecedence = 'secondary',
            updatedAt = NOW()
        WHERE contacts.id IN (email_primary_id, phone_primary_id)
        AND contacts.id != chosen_primary_id;

        UPDATE contacts
        SET linkedId = chosen_primary_id,
            updatedAt = NOW()
        WHERE contacts.linkedId IN (email_primary_id, phone_primary_id);

        RETURN QUERY
        SELECT 
            contacts.id, contacts.email, contacts.phoneNumber, contacts.linkedId,
            contacts.linkPrecedence, contacts.createdAt, contacts.updatedAt, contacts.deletedAt
        FROM contacts
        WHERE contacts.id = chosen_primary_id OR contacts.linkedId = chosen_primary_id
        ORDER BY contacts.createdAt;
        RETURN;
    END IF;

    IF email_primary_id IS NOT NULL THEN
        INSERT INTO contacts (email, phoneNumber, linkedId, linkPrecedence)
        VALUES (p_email, p_phone, email_primary_id, 'secondary');

        RETURN QUERY
        SELECT 
            contacts.id, contacts.email, contacts.phoneNumber, contacts.linkedId,
            contacts.linkPrecedence, contacts.createdAt, contacts.updatedAt, contacts.deletedAt
        FROM contacts
        WHERE contacts.id = email_primary_id OR contacts.linkedId = email_primary_id
        ORDER BY contacts.createdAt;
        RETURN;
    END IF;

    IF phone_primary_id IS NOT NULL THEN
        INSERT INTO contacts (email, phoneNumber, linkedId, linkPrecedence)
        VALUES (p_email, p_phone, phone_primary_id, 'secondary');

        RETURN QUERY
        SELECT 
            contacts.id, contacts.email, contacts.phoneNumber, contacts.linkedId,
            contacts.linkPrecedence, contacts.createdAt, contacts.updatedAt, contacts.deletedAt
        FROM contacts
        WHERE contacts.id = phone_primary_id OR contacts.linkedId = phone_primary_id
        ORDER BY contacts.createdAt;
        RETURN;
    END IF;

    RETURN;
END;
$$;`

const ADD_TEST_DATA = `INSERT INTO contacts (id, email, phoneNumber, linkedId, linkPrecedence) 
	VALUES 
	(1, 'email1', 'phone1', NULL, 'primary'),
	(2, 'email2', 'phone2', NULL,'primary'),
	(3, 'email3', 'phone2', 2, 'secondary'),
	(4, 'email2', 'phone4', 2,'secondary'),
	(5, 'email5', 'phone1', 1, 'secondary');
`
