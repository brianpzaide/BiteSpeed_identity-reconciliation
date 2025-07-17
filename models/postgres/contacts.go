package postgres

const CREATE_CONTACTS_TABLE = `CREATE TABLE IF NOT EXISTS contacts (
    id SERIAL PRIMARY KEY,
    phoneNumber TEXT,
    email TEXT,
    linkedId INTEGER REFERENCES contacts(id),
    linkPrecedence TEXT CHECK (linkPrecedence IN ('primary', 'secondary')),
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP
)`

const CREATE_STORED_PROCEDURE_WITH_ADVISORY_LOCK = `CREATE OR REPLACE FUNCTION insert_contact(p_email TEXT, p_phone TEXT)
	RETURNS contacts AS $$
		DECLARE
		    existing_primary_id INT;
			existing_email_match INT;
			existing_phone_number_match INT;
		    new_id INT;
		BEGIN
	    	PERFORM pg_advisory_xact_lock(
	    	    COALESCE(hashtext(p_email), 0)
	    	    + COALESCE(hashtext(p_phone), 0)
	    	);

	    	SELECT id
	    	INTO existing_primary_id
	    	FROM contacts
	    	WHERE (phoneNumber = p_phone OR email = p_email)
	    	  AND linkPrecedence = 'primary'
	    	LIMIT 1;

	    	IF existing_primary_id IS NULL THEN
	    	    INSERT INTO contacts (
	    	        phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt
	    	    )
	    	    VALUES (
	    	        p_phone, p_email, NULL, 'primary', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	    	    );
	    	ELSE
				SELECT 1
	    		INTO existing_phone_number_match
	    		FROM contacts
	    		WHERE phoneNumber = p_phone AND (id = existing_primary_id OR linkedId = existing_primary_id)
	    		LIMIT 1;

				SELECT 1
	    		INTO existing_email_match
	    		FROM contacts
	    		WHERE email = p_email AND (id = existing_primary_id OR linkedId = existing_primary_id)
	    		LIMIT 1;			
				
				IF existing_email_match IS NULL OR existing_phone_number_match IS NULL THEN
	        		INSERT INTO contacts (
	        		    phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt
	        		)
	        		VALUES (
	        		    p_phone, p_email, existing_primary_id, 'secondary', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	        		);
	    		END IF;
			END IF;

	    	RETURN QUERY
    		SELECT * FROM contacts
    		WHERE id = existing_primary_id OR linkedId = existing_primary_id;
	END;
	$$ LANGUAGE plpgsql;`
