package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initEnv()

	// check db connection
	println("Checking database connection...")
	checkDB()
	println("Database connection success")

	serverExit := &sync.WaitGroup{}
	serverExit.Add(1)

	// start server
	println("Starting server...")
	server := startServer(serverExit)
	println("Server started at", os.Getenv("PDS_PORT"))

	// if enter pressed
	println("\nPress enter key to shut down")
	fmt.Scanln()

	// shutdown
	println("Server shutting down...")
	err := server.Shutdown(context.TODO())
	if err != nil {
		panic(err)
	}

	// wait server to exit
	serverExit.Wait()
	println("Done exiting")
}
