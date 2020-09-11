package main

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

func main() {
	// initEnv()

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
	println("\nmain: Type exit then enter to shut down")
	println("main: Server started at", server.Addr, "if no error occured")

	exit := ""
	for exit != "exit" {
		fmt.Scanln(&exit)
	}

	println("\nmain: Server shutting down...")
	sessionCleanerExit <- 0
	err := server.Shutdown(context.TODO())
	if err != nil {
		panic(err)
	}
	// wait server to exit
	serverExit.Wait()
	println("main: Done exiting")
}
