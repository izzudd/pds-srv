package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
)

type authData struct {
	Uname string `json:"username"`
	Pass  string `json:"password"`
}

type userData struct {
	ID      string `json:"uid"`
	Name    string `json:"name"`
	Session string `json:"sid"`
}

type voteData struct {
	ID      string `json:"uid"`
	Session string `json:"sid"`
	Value   int    `json:"value"`
}

type outData struct {
	ID      string `json:"uid"`
	Session string `json:"sid"`
}

func parseBody(r io.Reader, v interface{}) error {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func randomString(length int) string {
	const s = "qwertyuiopasdfghjklzxcvbnm1234567890QWERTYUIOPASDFGHJKLZXCVBNM"
	var seed *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	res := make([]byte, length)
	for i := range res {
		res[i] = s[seed.Intn(len(s))]
	}
	return string(res)
}
