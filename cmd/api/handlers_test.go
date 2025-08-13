package main

import (
	"bitespeed_task/models/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setup(t *testing.T, dbType string) *application {
	testApp := &application{}
	switch {
	case dbType == "postgres":
		m, err := store.New("postgres", POSTGRES_DSN)
		if err != nil {
			t.Fatal(err)
		}
		testApp.models = m
	default:
		m, err := store.New("sqlite", SQLITE_DSN)
		if err != nil {
			t.Fatal(err)
		}
		testApp.models = m
	}
	return testApp
}

var expectedResults = []*ReconciliationResponse{
	{PrimaryContactId: 1,
		Emails:              []string{"email1"},
		PhoneNumbers:        []string{"phone1"},
		SecondaryContactIds: []int64{},
	},

	{PrimaryContactId: 2,
		Emails:              []string{"email2"},
		PhoneNumbers:        []string{"phone2"},
		SecondaryContactIds: []int64{},
	},

	{PrimaryContactId: 2,
		Emails:              []string{"email2", "email3"},
		PhoneNumbers:        []string{"phone2"},
		SecondaryContactIds: []int64{3},
	},

	{PrimaryContactId: 2,
		Emails:              []string{"email2", "email3"},
		PhoneNumbers:        []string{"phone2", "phone4"},
		SecondaryContactIds: []int64{3, 4},
	},

	{PrimaryContactId: 1,
		Emails:              []string{"email1", "email5"},
		PhoneNumbers:        []string{"phone1"},
		SecondaryContactIds: []int64{5},
	},

	{PrimaryContactId: 1,
		Emails:              []string{"email1", "email2", "email3", "email5"},
		PhoneNumbers:        []string{"phone1", "phone2", "phone4"},
		SecondaryContactIds: []int64{2, 3, 4, 5},
	},
}

func Test_reconciliate(t *testing.T) {
	testApp := setup(t, "postgres")

	tests := []struct {
		name                   string
		requestBody            string
		expectedResultObjectId int
	}{
		{"primary insert 1", `{"email":"email1","phoneNumber":"phone1"}`, 0},
		{"primary insert 2", `{"email":"email2","phoneNumber":"phone2"}`, 1},
		{"secondary insert to 2 new email", `{"email":"email3","phoneNumber":"phone2"}`, 2},
		{"secondary insert to 2 new phone", `{"email":"email2","phoneNumber":"phone4"}`, 3},
		{"secondary insert to 1 new email", `{"email":"email5","phoneNumber":"phone1"}`, 4},
		{"secondary non insert nothing new", `{"email":"email5","phoneNumber":"phone4"}`, 5},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/identify", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(testApp.reconciliate)
			handler.ServeHTTP(rr, req)

			// if rr.Code != http.StatusOK {
			// 	t.Fatalf("unexpected status code: got %d, body: %s", rr.Code, rr.Body.String())
			// }

			// rawResp, _ := io.ReadAll([]byte(rr.Body.String()))
			// if len(rawResp) == 0 {
			// 	t.Fatalf("empty response body for test case %d", i)
			// }

			var returned struct {
				Contact ReconciliationResponse `json:"contact"`
			}
			rawResp := rr.Body.Bytes()
			// t.Log(string(rawResp))
			if err := json.Unmarshal(rawResp, &returned); err != nil {
				t.Fatalf("failed to decode JSON: %v\nRaw body: %s", err, string(rawResp))
			}

			if !compareObjects(expectedResults[tc.expectedResultObjectId], &returned.Contact) {
				t.Errorf("mismatch for test case %d: got %+v, want %+v", i, returned, expectedResults[tc.expectedResultObjectId])
			}
		})
	}
}

func compareObjects(expected, returned *ReconciliationResponse) bool {
	pIdMatch := expected.PrimaryContactId == returned.PrimaryContactId
	sIdsMatch := checkSlicesEqual(expected.SecondaryContactIds, returned.SecondaryContactIds)
	phoneNumbersMatch := checkSlicesEqual(expected.PhoneNumbers, returned.PhoneNumbers)
	emailsMatch := checkSlicesEqual(expected.Emails, returned.Emails)

	return pIdMatch && sIdsMatch && emailsMatch && phoneNumbersMatch
}

func checkSlicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	seen := make(map[T]bool, len(a))
	for _, v := range a {
		seen[v] = true
	}

	for _, v := range b {
		if !seen[v] {
			return false
		}
	}
	return true
}
