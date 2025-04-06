package security

import (
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"stout.dev/login/internal/pkg/postgres"
)

func CreateSession(password string, account postgres.Account, err error, w http.ResponseWriter, conn *pgx.Conn) bool {
	if passed := CheckPasswordHash(password, account.Password); !passed {
		log.Printf("Error checking password: %v", err)
		http.Error(w, UnauthorizedMessage, http.StatusUnauthorized)
		return true
	}

	sessionToken := GenerateSecureToken()
	tokenExpiration := time.Now().Add(time.Hour * 1) // 1 hours

	if err = postgres.CreateSessionToken(conn, sessionToken, account.ID, tokenExpiration); err != nil {
		log.Printf("Error creating session token: %v", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: true,
	})

	csrfToken := GenerateSecureToken()

	if err = postgres.CreateCsrfToken(conn, csrfToken, account.ID, tokenExpiration); err != nil {
		log.Printf("Error creating csrf token: %v", err)
		http.Error(w, InternalServerErrorMessage, http.StatusInternalServerError)
		return true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: false,
	})
	return false
}
