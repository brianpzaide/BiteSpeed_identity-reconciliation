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
			primary_existing_email_match INT;
			existing_email_created_at TIMESTAMP;
			existing_phone_number_match INT;
			primary_existing_phone_number_match INT;
			existing_phone_number_created_at TIMESTAMP;
		    new_id INT;
		BEGIN
	    	PERFORM pg_advisory_xact_lock(
	    	    COALESCE(hashtext(p_email), 0), 
				COALESCE(hashtext(p_phone), 0)
	    	);

	    	SELECT id
	    	INTO existing_primary_id
	    	FROM contacts
	    	WHERE phoneNumber = p_phone OR email = p_email
			ORDER BY createdAt
	    	LIMIT 1;

	    	IF existing_primary_id IS NULL THEN
	    	    INSERT INTO contacts (
	    	        phoneNumber, email, linkPrecedence
	    	    )
	    	    VALUES (
	    	        p_phone, p_email, 'primary'
	    	    );
	    	ELSE
				SELECT id, linkedId, createdAt
	    		INTO existing_phone_number_match, primary_existing_phone_number_match, existing_phone_number_created_at
	    		FROM contacts
	    		WHERE phoneNumber = p_phone
				ORDER BY createdAT
	    		LIMIT 1;

				SELECT id, linkedId, createdAt
	    		INTO existing_email_match, primary_existing_email_match, existing_email_created_at
	    		FROM contacts
	    		WHERE email = p_email
				ORDER BY createdAT
	    		LIMIT 1;
				
				IF existing_email_match IS NULL AND existing_phone_number_match IS NOT NULL THEN
					IF primary_existing_phone_number_match IS NOT NULL THEN
	        			INSERT INTO contacts (
	        			    phoneNumber, email, linkedId, linkPrecedence
	        			)
	        			VALUES (
	        			    p_phone, p_email, primary_existing_phone_number_match, 'secondary'
	        			);
					ELSE
						INSERT INTO contacts (
	        			    phoneNumber, email, linkedId, linkPrecedence
	        			)
	        			VALUES (
	        			    p_phone, p_email, existing_phone_number_match, 'secondary'
	        			);
					END IF;	
	    		END IF;
				
				IF existing_phone_number_match IS NULL AND existing_email_match IS NOT NULL THEN
					IF primary_existing_email_match IS NOT NULL THEN
	        			INSERT INTO contacts (
	        			    phoneNumber, email, linkedId, linkPrecedence
	        			)
	        			VALUES (
	        			    p_phone, p_email, primary_existing_email_match, 'secondary'
	        			);
					ELSE
						INSERT INTO contacts (
	        			    phoneNumber, email, linkedId, linkPrecedence
	        			)
	        			VALUES (
	        			    p_phone, p_email, existing_email_match, 'secondary'
	        			);
					END IF;
	    		END IF;

				IF (existing_email_match IS NOT NULL) AND (existing_phone_number_match IS NOT NULL) THEN
					IF (primary_existing_email_match IS NOT NULL) AND (primary_existing_phone_number_match IS NOT NULL) AND (primary_existing_phone_number_match != primary_existing_email_match) THEN
						IF existing_phone_number_created_at < existing_email_created_at THEN
	        				UPDATE contacts
							SET linkedId = primary_existing_phone_number_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        				WHERE id = existing_email_match OR linkedId = existing_email_match;
						ELSE
	        				UPDATE contacts
							SET linkedId = existing_email_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        				WHERE id = existing_phone_number_match OR linkedId = existing_phone_number_match;		
						END IF;				
					END IF;

					IF (primary_existing_email_match IS NULL) AND (primary_existing_phone_number_match IS NOT NULL) AND (primary_existing_phone_number_match != existing_email_match) THEN
						IF phone_number_created_at < primary_email_created_at THEN
	        				UPDATE contacts
							SET linkedId = primary_existing_phone_number_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        				WHERE id = existing_email_match OR linkedId = existing_email_match;
						ELSE
	        				UPDATE contacts
							SET linkedId = existing_email_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        				WHERE id = existing_phone_number_match OR linkedId = existing_phone_number_match;		
						END IF;				
					END IF;
					
					
					IF phone_number_created_at < email_created_at THEN
	        			UPDATE contacts
						SET linkedId = existing_phone_number_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        			WHERE id = existing_email_match OR linkedId = existing_email_match;
					ELSE
	        			UPDATE contacts
						SET linkedId = existing_email_match, linkPrecedence = 'secondary', updatedAt = NOW()
	        			WHERE id = existing_phone_number_match OR linkedId = existing_phone_number_match;		
					END IF;				
	    		END IF;
			END IF;

	    	RETURN QUERY
			WITH primary_id AS (SELECT id FROM contacts WHERE (email = p_email OR phoneNumber = p_phone) AND linkPrecedence = 'primary')
    		SELECT * FROM contacts
    		WHERE id IN (SELECT * FROM primary_id) 
			OR linkedId IN (SELECT * FROM primary_id)
			ORDER BY createdAt;
	END;
	$$ LANGUAGE plpgsql;`

const ADD_TEST_DATA = `INSERT INTO contacts (phoneNumber, email, linkedId, linkPrecedence) 
	VALUES 
	('phone1', 'email1', NULL, 'primary'),
	('phone2', 'email2', NULL,'primary');
`
