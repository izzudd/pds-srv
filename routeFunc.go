package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

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
	var isVoted bool
	err = db.QueryRow("SELECT id, name, done FROM public.datasiswa WHERE username = $1 AND password = $2", auth.Uname, auth.Pass).Scan(&data.ID, &data.Name, &isVoted)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	// user voted
	if isVoted {
		http.Error(w, "user voted", 400)
		return
	}

	// delete session if exist
	_, err = db.Exec("DELETE FROM session WHERE user_id = $1", data.ID)
	if err != nil {
		http.Error(w, "session exist cant drop: "+err.Error(), 400)
		return
	}

	// apply session
	data.Session = randomString(16)
	_, err = db.Exec("INSERT INTO session (session_id, user_id, expired) VALUES ($1, $2, CURRENT_TIMESTAMP + INTERVAL '15 MINUTE')", data.Session, data.ID)
	if err != nil {
		http.Error(w, "cant assign session: "+err.Error(), 400)
		return
	}

	// write response if success
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to send response: "+err.Error(), 500)
		// if can't send response delete session drom table
		_, err = db.Exec("DELETE FROM session WHERE session_id = $1", data.Session)
		if err != nil {
			http.Error(w, err.Error(), 400)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
	err = db.QueryRow("SELECT count(*) FROM session WHERE session_id = $1 AND user_id = $2", vote.Session, vote.ID).Scan(&valid)
	if valid != 1 {
		http.Error(w, "invalid session", 400)
		return
	}

	// check if user voted
	var isVoted bool
	err = db.QueryRow("SELECT done FROM datasiswa WHERE id = $1", vote.ID).Scan(&isVoted)
	if isVoted {
		http.Error(w, "user voted", 400)
		return
	}

	// apply vote data to db
	_, err = db.Exec("UPDATE result SET value = value + 1 WHERE id = $1", vote.Value)
	if err != nil {
		http.Error(w, "vote data not applied: "+err.Error(), 400)
		return
	}
	// set user status to voted
	// _, err = db.Exec("UPDATE datasiswa SET done = true WHERE id = $1", vote.ID)
	// if err != nil {
	// 	http.Error(w, "vote data not applied: "+err.Error(), 400)
	// 	return
	// }

	// send success response
	w.WriteHeader(200)
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
	_, err = db.Exec("DELETE FROM session WHERE session_id = $1 AND user_id = $2", out.Session, out.ID)
	if err != nil {
		http.Error(w, "failed to logout: "+err.Error(), 400)
		return
	}

	// send response if success
	w.WriteHeader(200)
}
