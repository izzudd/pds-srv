package main

import (
	"database/sql"
	"os"
)

// DB is the databse singleton
var DB *sql.DB

// InitDB initialize database to server
func InitDB() {
	println("db: Checking database connection...")
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	println("db: Database connection success")
}
