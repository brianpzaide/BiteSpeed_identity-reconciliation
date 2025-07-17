package main

import (
	"bitespeed_task/models"
	"net/http"
)

type ReconciliationResponse struct {
	PrimaryContatctId   int64    `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNmbers"`
	SecondaryContactIds []int64  `json:"secondaryContactIds"`
}

func registerUser(app *application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := &models.QueryREquest{}
		err := app.readJSON(w, r, input)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		contacts, err := app.models.Contacts.Fetch(input)

		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		err = app.writeJSON(w, http.StatusCreated, envelope{"contact": contacts}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
