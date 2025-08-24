# BiteSpeed_identity-reconciliation
This repository contains my solution to the [technical task](https://bitespeed.notion.site/Bitespeed-Backend-Task-Identity-Reconciliation-1fb21bb2a930802eb896d4409460375c) for the [Backend Developer - SDE 1](https://bitespeed.notion.site/Backend-Developer-SDE-1-357cd0ddceba497bbf5f4dc88b03522b) position at [BiteSpeed](https://www.bitespeed.co/).

## 1. Problem Statement

  Businesses often struggle with fragmented customer contact data. A single customer may appear multiple times in the database with different combinations of email addresses and phone numbers. For example, one record might have only the phone number, while another has only the email address. Without reconciliation, this leads to duplication, inconsistency, and difficulty in retrieving a customer’s complete information.

  The challenge is to **reconcile contacts automatically** so that all data points (emails and phone numbers) belonging to the same customer can be linked together under a primary contact.

  Design a web service with an endpoint ```/identify``` that will receive HTTP POST requests with JSON body of the following format:

  ```typescript
  {
	  "email"?: string,
	  "phoneNumber"?: string
  }
  ```
  The web service should return an HTTP 200 response with a JSON payload containing the consolidated contact.
  Your response should be in this format:

  ```typescript
	{
		"contact":{
			"primaryContatctId": number,
			"emails": string[], // first element being email of primary contact 
			"phoneNumbers": string[], // first element being phoneNumber of primary contact
			"secondaryContactIds": number[] // Array of all Contact IDs that are "secondary" to the primary contact
		}
	}

  ```
#### Database schema
``` typescript 
  id                   Int                   
  phoneNumber          String?
  email                String?
  linkedId             Int? // the ID of another Contact linked to this one
  linkPrecedence       "secondary"|"primary" // "primary" if it's the first Contact in the link
  createdAt            DateTime              
  updatedAt            DateTime              
  deletedAt            DateTime?
```
* Customers can have multiple **`Contact`** rows in the database against them. All of the rows are linked together with the oldest one being treated as "primary” and the rest as “secondary”.
* **`Contact`** rows are linked if they have either of **`email`** or **`phone`** as common.

## Example Workflow

#### Start with empty database.
➡️ Request 1:
```json
{
	"email": "email1",
	"phoneNumber": "phone1"
}
```
Creates a primary record:

```table
id |   email   | phoneNumber | linkedId | linkPrecedence | createdAt
---+-----------+-------------+----------+----------------+------------
 1 | email1    | phone1      | NULL     | primary        | 2023-01-01
```

➡️ Request 2:

 ```json
{
	"email": "email2",
	"phoneNumber": "phone1"
}
```
Matches existing phone number → insert as secondary:

```table
| id | email  | phoneNumber | linkedId | linkPrecedence | createdAt  |
| -- | ------ | ----------- | -------- | -------------- | ---------- |
| 1  | email1 | phone1      | NULL     | primary        | 2023-01-01 |
| 2  | email2 | phone1      | 1        | secondary      | 2023-01-02 |
```

### Summary:

```table
| Case                                         | Action                                              |
| -------------------------------------------- | --------------------------------------------------- |
| Neither email nor phone exists               | Insert as `primary`                                 |
| Only one (email or phone) exists             | Insert as `secondary` linked to the match           |
| Both exist, and match same primary           | Do nothing                                          |
| Both exist, but point to different primaries | Link them by demoting the newer primary — no insert |
```
---

## 2. Solution Features

* **Repository Pattern**
  
  This helps in keeping the application layer clean and decoupled from the database implementation. Currently supporting two database systems **PostgreSQL** and **SQLite3**.

* **Database Integrity & Concurrency Handling**

  * **PostgreSQL:** For PostgreSQL the solution uses a stored procedure with advisory locks to prevent race conditions and avoid data corruption during concurrent writes.
  * **SQLite3:**  For SQLite3 the solution wraps reconciliation operations in **transactions** to ensure atomicity and consistency.

* **Automatic Reconciliation**
  
  New contacts are automatically reconciled with existing records if overlaps exist (same email or phone number).

* **Minimized Tree Structure**
  
  All secondary contacts are linked directly to a **root primary contact**. This avoids complex recursive queries when retrieving all customer data.

---

## 3. Notes & Limitations

* **Empty Database Requirement**
  
  This solution assumes the database starts empty. If you already have existing data, reconciliation correctness cannot be guaranteed.
  Recommended usage: start with a fresh database and add records incrementally. Each new record will be automatically reconciled at insert time.

* **Reconciliation Strategy**
  
  Inorder to ensure **flat structure** and avoid recursive queries for fetching related records all child(secondary) contacts are attached to the root primary contact.


