package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"stout.dev/login/postgres"
)

const (
	internalServerErrorMessage = "Banana machine not working"
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
		handleLogout(w, r, conn)
	})

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func handleLogout(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {

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

	sessionToken := generateSecureToken()
	tokenExpiration := time.Now().Add(time.Hour * 1) // 1 hours

	if err = postgres.CreateSessionToken(conn, username, sessionToken, tokenExpiration); err != nil {
		log.Printf("Error creating session token: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  tokenExpiration,
		HttpOnly: true,
	})

	csrfToken := generateSecureToken()

	if err = postgres.CreateCsrfToken(conn, username, sessionToken, tokenExpiration); err != nil {
		log.Printf("Error creating csrf token: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  tokenExpiration,
		HttpOnly: false,
	})

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
		http.Error(w, "", http.StatusForbidden)
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
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	if err = postgres.CreateAccount(conn, username, hashedPassword); err != nil {
		log.Printf("Error creating account: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Minion registered successfully"))
}
