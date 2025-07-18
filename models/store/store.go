package store

import (
	"bitespeed_task/models"
	"bitespeed_task/models/store/postgres"
)

const SQLITE_DSN = "./identity_reconciliation.db"

type ContactModelInterface interface {
	Reconciliate(email, phoneNumber string) ([]*models.Contact, error)
	Close()
}

type Models struct {
	ContactsModel ContactModelInterface
}

func New(dbType, dsn string) (*Models, error) {
	switch dbType {
	case "postgres":
		cm, err := postgres.NewPostgresModel(dsn)
		if err != nil {
			return nil, err
		}
		return &Models{ContactsModel: cm}, nil
	default:
		cm, err := postgres.NewPostgresModel(dsn)
		if err != nil {
			return nil, err
		}
		return &Models{ContactsModel: cm}, nil
	}
}
