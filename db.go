package main

import (
	"database/sql"
	"os"
)

// ConnectDB begin connection to database
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("PDS_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func checkDB() {
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	db.Close()
}
