package v1

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"stout.dev/login/internal/pkg/postgres"
	"stout.dev/login/internal/pkg/security"
)

const (
	internalServerErrorMessage = "Minion broke the banana machine not working"
	methodNotAllowedMessage    = "Minion not allowed"
	statusUnauthorizedMessage  = "Minion is unauthorized"
	statusForbiddenMessage     = "Minion tried something naughty"
)

type LoginV1 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleLogout(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	if r.Method != http.MethodGet {
		http.Error(w, methodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, statusForbiddenMessage, http.StatusForbidden)
		return
	}

	csrfToken := r.Header.Get("x-csrf-token")
	if csrfToken == "" {
		http.Error(w, statusForbiddenMessage, http.StatusForbidden)
		return
	}

	accountID, err := postgres.GetAccountIdBySessionToken(conn, sessionToken.Value)
	if err != nil {
		log.Printf("Error getting account id: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	csrfAccountID, err := postgres.GetCsrfIdBySessionToken(conn, csrfToken)
	if err != nil {
		log.Printf("Error getting csrf account id: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	if accountID != csrfAccountID {
		http.Error(w, statusForbiddenMessage, http.StatusForbidden)
		return
	}

	logoutAllDevices := r.URL.Query().Get("all_devices") == "true"
	if logoutAllDevices {
		if err := postgres.DeleteCsrfTokensByAccountId(conn, accountID); err != nil {
			log.Printf("Error deleting csrf tokens: %v", err)
			http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		if err := postgres.DeleteSessionTokensByAccountId(conn, accountID); err != nil {
			log.Printf("Error deleting session tokens: %v", err)
			http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
			return
		}
	} else {
		if err := postgres.DeleteCsrfToken(conn, csrfToken); err != nil {
			log.Printf("Error deleting csrf token: %v", err)
			http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
			return
		}

		if err := postgres.DeleteSessionToken(conn, sessionToken.Value); err != nil {
			log.Printf("Error deleting session token: %v", err)
			http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion logged out successfully"))
}

func HandleLogin(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	var loginData LoginV1

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	username := loginData.Username
	password := loginData.Password

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	account, err := postgres.GetAccount(conn, username)
	if err != nil {
		log.Printf("Error getting password: %v", err)
		http.Error(w, statusUnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	if passed := security.CheckPasswordHash(password, account.Password); !passed {
		log.Printf("Error checking password: %v", err)
		http.Error(w, statusUnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	sessionToken := security.GenerateSecureToken()
	tokenExpiration := time.Now().Add(time.Hour * 1) // 1 hours

	if err = postgres.CreateSessionToken(conn, sessionToken, account.ID, tokenExpiration); err != nil {
		log.Printf("Error creating session token: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: true,
	})

	csrfToken := security.GenerateSecureToken()

	if err = postgres.CreateCsrfToken(conn, csrfToken, account.ID, tokenExpiration); err != nil {
		log.Printf("Error creating csrf token: %v", err)
		http.Error(w, internalServerErrorMessage, http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: false,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion logged in successfully"))
}

func HandleRegister(w http.ResponseWriter, r *http.Request, conn *pgx.Conn) {
	if r.Method != http.MethodPost {
		http.Error(w, methodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	var loginData LoginV1

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	username := loginData.Username
	password := loginData.Password

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
	hashedPassword, err := security.HashPassword(password)
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
