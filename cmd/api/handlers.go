package main

import (
	"bitespeed_task/models"
	"fmt"
	"net/http"
)

type QueryRequest struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

type ReconciliationResponse struct {
	PrimaryContactId    int64    `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNmbers"`
	SecondaryContactIds []int64  `json:"secondaryContactIds"`
}

func (app *application) reconciliate(w http.ResponseWriter, r *http.Request) {
	input := &QueryRequest{}
	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	contacts, err := app.models.Reconciliate(input.Email, input.PhoneNumber)
	fmt.Printf("email: %s, phoneNumber: %s\n", input.Email, input.PhoneNumber)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := generateresponse(contacts)

	err = app.writeJSON(w, http.StatusCreated, envelope{"contact": data}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func generateresponse(contacts []*models.Contact) *ReconciliationResponse {

	var primaryContactId int64 = 0
	secondaryContactIds := make([]int64, 0)
	emailsMap := make(map[string]bool)
	phonesMap := make(map[string]bool)
	emails := make([]string, 0)
	phoneNumbers := make([]string, 0)

	for i, contact := range contacts {
		if i == 0 {
			primaryContactId = contact.ID
		} else {
			secondaryContactIds = append(secondaryContactIds, contact.ID)
		}
		if ok := emailsMap[contact.Email]; !ok {
			emails = append(emails, contact.Email)
			emailsMap[contact.Email] = true
		}

		if ok := phonesMap[contact.PhoneNumber]; !ok {
			phoneNumbers = append(phoneNumbers, contact.PhoneNumber)
			phonesMap[contact.PhoneNumber] = true
		}
	}

	return &ReconciliationResponse{
		PrimaryContactId:    primaryContactId,
		SecondaryContactIds: secondaryContactIds,
		Emails:              emails,
		PhoneNumbers:        phoneNumbers,
	}
}
