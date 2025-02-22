package app

import (
	"log"
	"net/http"

	v1 "stout.dev/jisho/internal/app/v1"
)

func NewApp() {
	http.HandleFunc("/v1/search", v1.HandleJishoReponse)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
