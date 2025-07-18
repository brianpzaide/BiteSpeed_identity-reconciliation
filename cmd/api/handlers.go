package main

import (
	"net/http"
)

type QueryRequest struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNnumber,omitempty"`
}

type ReconciliationResponse struct {
	PrimaryContatctId   int64    `json:"primaryContatctId"`
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

	contacts, err := app.models.ContactsModel.Reconciliate(input.Email, input.PhoneNumber)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"contact": contacts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
