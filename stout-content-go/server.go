package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const CV_PATH = "./assets/Teun-Johán-Stout.pdf"

func main() {
	http.HandleFunc("/v1/content/cv", handleDownloadCv)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleDownloadCv(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(CV_PATH)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	defer f.Close()

	//Set header
	w.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}
