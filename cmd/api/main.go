package main

import (
	"bitespeed_task/models/store"
	"flag"
	"log"
	"os"
)

const SQLITE_DSN = "./identity_reconciliation.db?_txlock=immediate"

const POSTGRES_DSN = "host=localhost port=5432 user=postgres password=mysecretpassword dbname=identityreconciliation sslmode=disable timezone=UTC connect_timeout=5"

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models store.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.db.dsn, "db-dsn", POSTGRES_DSN, "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var (
		m   store.Models
		err error
	)

	if cfg.db.dsn == "" {
		m, err = store.New("sqlite", SQLITE_DSN)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		m, err = store.New("postgres", cfg.db.dsn)
		if err != nil {
			logger.Fatal(err)
		}
	}

	defer m.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: m,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
