package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	initEnv()

	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/vote", vote).Methods("POST")
	r.HandleFunc("/logout", logout).Methods("POST")

	// check db connection
	db, err := ConnectDB()
	if err != nil {
		panic(err)
	}
	db.Close()

	server := http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	println("Server Listen at", os.Getenv("PDS_PORT"))
	log.Fatal((&server).ListenAndServe())
}

// http method function

func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var auth authData
	err := parseBody(r.Body, &auth)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// making data in db
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	var data userData
	err = db.QueryRow("SELECT id, name FROM datasiswa WHERE username = ? AND password = ?", auth.Uname, auth.Pass).Scan(&data.ID, &data.Name)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// delete session if exist
	_, err = db.Exec("DELETE FROM session WHERE user_id = ?", data.ID)
	if err != nil {
		http.Error(w, "session exist cant drop: "+err.Error(), 400)
		return
	}

	// apply session
	data.Session = randomString(16)
	_, err = db.Exec("INSERT INTO session (`session_id`, `user_id`) VALUES (?, ?)", data.Session, data.ID)
	if err != nil {
		http.Error(w, "cant assign session: "+err.Error(), 400)
		return
	}

	// write response if success
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to send response: "+err.Error(), 500)
		// if can't send response delete session drom table
		_, err = db.Exec("DELETE FROM session WHERE session_id = ?", data.Session)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		return
	}
	w.Write(res)
}

func vote(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var vote voteData
	err := parseBody(r.Body, &vote)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// init db to communicate
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	// check session validity
	var valid int
	err = db.QueryRow("SELECT count(*) FROM session WHERE session_id = ? AND user_id = ?", vote.Session, vote.ID).Scan(&valid)
	// err = db.QueryRow("SELECT count(*) FROM session WHERE session_id = 'ZRTemKMFbdOgFqhe' AND user_id = '1'").Scan(&valid)
	if valid != 1 {
		w.Write([]byte("invalid session"))
		return
	}

	// check if user voted
	var isVoted bool
	err = db.QueryRow("SELECT done FROM datasiswa WHERE id = ?", vote.ID).Scan(&isVoted)
	if isVoted {
		w.Write([]byte("user voted"))
		return
	}

	// apply vote data to db
	_, err = db.Exec("UPDATE result SET value = value + 1 WHERE id = ?", vote.Value)
	if err != nil {
		http.Error(w, "vote data not applied: "+err.Error(), 400)
		return
	}
	_, err = db.Exec("UPDATE datasiswa SET done = true WHERE id = ?", vote.ID)
	if err != nil {
		http.Error(w, "vote data not applied: "+err.Error(), 400)
		return
	}

	// send success response
	w.Write([]byte("success"))
}

func logout(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var out outData
	err := parseBody(r.Body, &out)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// init db to communicate
	db, err := ConnectDB()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer db.Close()

	// delete session
	_, err = db.Exec("DELETE FROM session WHERE session_id = ? AND user_id = ?", out.Session, out.ID)
	if err != nil {
		http.Error(w, "failed to logout: "+err.Error(), 400)
		return
	}

	// send response if success
	w.Write([]byte("success"))
}
