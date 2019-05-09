package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func registerHandlers() {
	http.HandleFunc("/", twillioHandler)
}

func twillioHandler(w http.ResponseWriter, r *http.Request) {
	// catch all, try to get out all the data we need, if not
	// found we will simply log and return a bad request. Not
	// the best practice but good enough for this purpose.

	r.ParseForm()
	phone := r.FormValue("phone")
	if phone == "" {
		log.Println("no phone number was sent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := r.FormValue("urls")
	urls := []string{}
	if err := json.Unmarshal([]byte(url), &urls); err != nil {
		log.Println("could not decode json from urls: ", err, urls)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code := processRequest(phone, urls)
	w.WriteHeader(code)
}

func runServer() {
	err := http.ListenAndServe(fHTTP, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
