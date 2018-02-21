package main

import (
	"log"
	"net/http"
)

func ok(w http.ResponseWriter, r *http.Request) {
	log.Printf("requested path: %s\n", r.URL)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
