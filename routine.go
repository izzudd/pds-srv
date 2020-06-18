package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func startServer(wg *sync.WaitGroup) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/vote", vote).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("POST")

	server := &http.Server{
		Handler:      r,
		Addr:         ":" + os.Getenv("PDS_PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		defer wg.Done()
		println("Starting server...")
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}

func startSessionCleanup(wg *sync.WaitGroup, exit <-chan int) {
	defer wg.Done()
	println("Starting session cleaner...")
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	println("Session cleaner started")
	for {
		select {
		default:
			_, err = db.Exec("DELETE FROM session WHERE expired < CURRENT_TIMESTAMP")
			if err != nil {
				log.Fatal(err.Error())
			}
		case <-exit:
			println("Session cleaner shut down")
			return
		}
	}
}
