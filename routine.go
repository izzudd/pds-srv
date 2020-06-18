package main

import (
	"net/http"
	"os"
	"sync"
	"time"

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
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return server
}
