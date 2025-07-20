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

const create_stored_procedure_with_advisory_lock = `CREATE OR REPLACE FUNCTION reconciliate(p_email TEXT, p_phone TEXT)
	RETURNS SETOF contacts AS $$
		DECLARE
		    existing_primary_id INT;
			existing_email_match INT;
			email_created_at TIMESTAMP;
			existing_phone_number_match INT;
			phone_number_created_at TIMESTAMP;
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
			ORDER BY created_at
	    	LIMIT 1;

	    	IF existing_primary_id IS NULL THEN
	    	    INSERT INTO contacts (
	    	        phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt
	    	    )
	    	    VALUES (
	    	        p_phone, p_email, NULL, 'primary', NOW(), NOW()
	    	    );
	    	ELSE
				SELECT id, created_at
	    		INTO existing_phone_number_match, phone_number_created_at
	    		FROM contacts
	    		WHERE phoneNumber = p_phone AND linkPrecedence = 'primary'
	    		LIMIT 1;

				SELECT id, created_at
	    		INTO existing_email_match, email_created_at
	    		FROM contacts
	    		WHERE email = p_email AND linkPrecedence = 'primary'
	    		LIMIT 1;
				
				IF existing_email_match IS NULL AND existing_phone_number_match IS NOT NULL THEN
	        		INSERT INTO contacts (
	        		    phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt
	        		)
	        		VALUES (
	        		    p_phone, p_email, existing_phone_number_match, 'secondary', NOW(), NOW()
	        		);
	    		END IF;
				
				IF existing_phone_number_match IS NULL AND existing_email_match IS NOT NULL THEN
	        		INSERT INTO contacts (
	        		    phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt
	        		)
	        		VALUES (
	        		    p_phone, p_email, existing_email_match, 'secondary', NOW(), NOW()
	        		);
	    		END IF;

				IF (existing_email_match IS NOT NULL) AND (existing_phone_number_match IS NOT NULL) AND (existing_phone_number_match != existing_email_match) THEN
					IF phone_number_created_at < email_created_at THEN
	        			UPDATE contacts
						SET linkedId = existing_phone_number_match, linkPrecedence = 'secondary', updated_at = NOW()
	        			WHERE id = existing_email_match OR linkedId = existing_email_match;
					ELSE
	        			UPDATE contacts
						SET linkedId = existing_email_match, linkPrecedence = 'secondary', updated_at = NOW()
	        			WHERE id = existing_phone_number_match OR linkedId = existing_phone_number_match;		
					END IF;				
	    		END IF;
			END IF;

	    	RETURN QUERY
    		SELECT * FROM contacts
    		WHERE email = p_email or phoneNumber = p_phone ORDER BY created_at;
	END;
	$$ LANGUAGE plpgsql;`
