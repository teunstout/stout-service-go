package app

import (
	"log"
	"net/http"

	v1 "stout.dev/jisho/internal/app/v1"
)

func NewApp() {
	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("/v1/search", v1.HandleJishoResponse)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
