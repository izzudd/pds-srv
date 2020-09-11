package main

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func startServer(wg *sync.WaitGroup) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/vote", vote).Methods("POST")
	r.HandleFunc("/api/logout", logout).Methods("POST")

	spa := spaHandler{staticPath: "dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.PathPrefix("/api").Subrouter().MethodNotAllowedHandler = http.HandlerFunc(notAllowed)

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
		handlers.AllowedMethods([]string{"POST"}),
		handlers.AllowedOrigins([]string{os.Getenv("PSD_ALLOWED_ORIGIN")}),
		handlers.AllowCredentials(),
	)
	server := &http.Server{
		Handler:      cors(r),
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		defer wg.Done()
		println("routine: Starting server...")
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			panic(err)
		}
		println("routine: Server shut down")
	}()

	return server
}

func startSessionCleanup(wg *sync.WaitGroup, exit <-chan int) {
	defer wg.Done()
	println("routine: Starting session cleaner...")
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	println("routine: Session cleaner started")
	for {
		select {
		default:
			_, err = db.Exec("DELETE FROM session WHERE expired < CURRENT_TIMESTAMP")
			if err != nil {
				log.Fatal(err.Error())
			}
		case <-exit:
			println("routine: Session cleaner shut down")
			return
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(404)
	w.Write([]byte("Not Found"))
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
}
