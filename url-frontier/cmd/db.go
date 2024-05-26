package main

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func (app *Config) connectToDB(maxRetries int) *sql.DB {
	dsn := os.Getenv("DATABASE_URL")

	for i := 0; i < maxRetries; i++ {
		connection, err := openDB(dsn)
		if err != nil {
			log.Printf("%s. Retrying in 2 seconds...", err.Error())
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("connected to postgres!, %s", dsn)
			return connection
		}
	}

	return nil
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
