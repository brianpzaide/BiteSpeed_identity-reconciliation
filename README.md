# BiteSpeed_identity-reconciliation
This repository contains my solution to the [technical task](https://bitespeed.notion.site/Bitespeed-Backend-Task-Identity-Reconciliation-1fb21bb2a930802eb896d4409460375c) for the [Backend Developer - SDE 1](https://bitespeed.notion.site/Backend-Developer-SDE-1-357cd0ddceba497bbf5f4dc88b03522b) position at [BiteSpeed](https://www.bitespeed.co/).

## 1. Problem Statement

  Businesses often struggle with fragmented customer contact data. A single customer may appear multiple times in the database with different combinations of email addresses and phone numbers. For example, one record might have only the phone number, while another has only the email address. Without reconciliation, this leads to duplication, inconsistency, and difficulty in retrieving a customer’s complete information.

  The challenge is to **reconcile contacts automatically** so that all data points (emails and phone numbers) belonging to the same customer can be linked together under a primary contact.

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
  
  All secondary contacts are linked directly to the **root primary contact**. This avoids complex recursive queries when retrieving all customer data.

---

## 3. Notes & Limitations

* **Empty Database Requirement**
  
  This solution assumes the database starts empty. If you already have existing data, reconciliation correctness cannot be guaranteed.
  Recommended usage: start with a fresh database and add records incrementally. Each new record will be automatically reconciled at insert time.

* **Reconciliation Strategy**
  
  Inorder to ensure **flat structure** and avoid recursive queries for fetching related records all child(secondary) contacts are attached to the root primary contact.

---

## 4. API Usage

### Endpoint: **`POST /identify`**

Used to insert a new contact and trigger reconciliation.

#### Request Body

```json
{
  "email": "foo@example.com",
  "phoneNumber": "1234567890"
}
```

* At least one of `email` or `phoneNumber` must be provided.
* If both match an existing record, the new contact is reconciled with that customer.

#### Example Responses

**Case 1: First Contact Inserted**

```json
{
  "contact": {
    "primaryContactId": 1,
    "emails": ["foo@example.com"],
    "phoneNumbers": ["1234567890"],
    "secondaryContactIds": []
  }
}
```

**Case 2: New Contact With Matching Email**

```json
{
  "contact": {
    "primaryContactId": 1,
    "emails": ["foo@example.com"],
    "phoneNumbers": ["1234567890", "9999999999"],
    "secondaryContactIds": [2]
  }
}
```

