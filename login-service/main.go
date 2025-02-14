package main

import (
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"stout.dev/login/postgres"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	http.HandleFunc("/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		// handleAuthenticate(w, r, conn)
	})

	http.HandleFunc("/v1/logout", func(w http.ResponseWriter, r *http.Request) {
		// handleLogout(w, r, conn)
	})

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleLogin(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := postgres.GetAccountPassword(conn, username)
	if err != nil {
		log.Printf("Error getting password: %v", err)
		http.Error(w, "Minion not found", http.StatusUnauthorized)
		return
	}

	if passed := checkPasswordHash(password, hashedPassword); !passed {
		log.Printf("Error checking password: %v", err)
		http.Error(w, "Minion not found", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion logged in successfully"))
}

func handleRegister(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	exists, err := postgres.AccountExists(conn, username)
	if err != nil {
		log.Printf("Error checking if account exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Minion already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err = postgres.CreateAccount(conn, username, hashedPassword); err != nil {
		log.Printf("Error creating account: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Minion registered successfully"))
}
