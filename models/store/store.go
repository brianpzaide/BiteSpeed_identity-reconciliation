package store

import (
	"bitespeed_task/models"
	"bitespeed_task/models/store/postgres"
	"bitespeed_task/models/store/sqlite"
)

const SQLITE_DSN = "./identity_reconciliation.db"

type ContactModelInterface interface {
	Reconciliate(email, phoneNumber string) ([]*models.Contact, error)
}

type Models interface {
	ContactModelInterface
	Close()
}

func New(dbType, dsn string) (Models, error) {
	switch dbType {
	case "postgres":
		cm, err := postgres.NewPostgresModel(dsn)
		if err != nil {
			return nil, err
		}
		return cm, nil
	default:
		cm, err := sqlite.NewSqliteModel(SQLITE_DSN)
		if err != nil {
			return nil, err
		}
		return cm, nil
	}
}
