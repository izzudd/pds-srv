package main

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initEnv()

	// check db connection
	checkDB()

	serverExit := &sync.WaitGroup{}
	serverExit.Add(2)
	sessionCleanerExit := make(chan int)

	// start server
	server := startServer(serverExit)

	// start session cleaner
	go startSessionCleanup(serverExit, sessionCleanerExit)

	// if enter pressed
	println("\nPress enter key to shut down")
	println("Server started at", server.Addr, "if no error occured")
	fmt.Scanln()

	// shutdown
	println("Server shutting down...")
	sessionCleanerExit <- 0
	err := server.Shutdown(context.TODO())
	if err != nil {
		panic(err)
	}

	// wait server to exit
	serverExit.Wait()
	println("Done exiting")
}
