package main

import (
	"log"
	"net/http"

	"stout.dev/login/postgres"
)

func main() {
	conn, err := postgres.Connect()
	if err != nil {
		panic(err)
	}

	if err := postgres.InitDatabase(conn); err != nil {
		panic(err)
	}

	http.HandleFunc("/v1/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, conn)
	})

	http.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, conn)
	})

	http.HandleFunc("/v1/logout", func(w http.ResponseWriter, r *http.Request) {
		handleLogout(w, r, conn)
	})

	// http.HandleFunc("/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
	// 	// handleAuthenticate(w, r, conn)
	// })

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
