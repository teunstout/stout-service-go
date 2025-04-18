package app

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	v1 "stout.dev/login/internal/app/v1"
	"stout.dev/login/internal/pkg/postgres"
	"stout.dev/login/internal/pkg/security"
)

var (
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	serverPort int
)

func NewApp() {
	conn, err := postgres.Connect()
	fatal(err)
	defer conn.Close(context.Background())

	privKeyPath := "/app/keys/app.rsa"
	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		fmt.Println("Private key not found, using default path")
		privKeyPath = `C:\Users\Teuns\Documents\Github\stout-service-go\app.rsa`
	}
	signBytes, err := os.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	pubKeyPath := "/app/keys/app.rsa.pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		fmt.Println("Public key not found, using default path")
		pubKeyPath = `C:\Users\Teuns\Documents\Github\stout-service-go\app.rsa.pub`
	}
	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	http.HandleFunc("/v1/register", func(w http.ResponseWriter, r *http.Request) {
		v1.HandleRegister(w, r, conn)
	})

	http.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
		v1.HandleLogin(w, r, conn, signKey)
	})

	http.HandleFunc("/v1/logout", func(w http.ResponseWriter, r *http.Request) {
		v1.HandleLogout(w, r, conn)
	})

	http.HandleFunc("/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		v1.HandleAuthenticate(w, r, conn)
	})

	http.HandleFunc("/v1/jwt", func(w http.ResponseWriter, r *http.Request) {
		security.CreateJwt(w, r, 1234, signKey)
	})

	http.HandleFunc("/v1/endpoints/restricted", restrictedHandler)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// only accessible with a valid token
func restrictedHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from request
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	}, request.WithClaims(&jwt.MapClaims{}))

	// If the token is missing or invalid, return error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token:", err)
		return
	}

	// Token is valid
	fmt.Fprintln(w, "Welcome,", token.Claims.(*jwt.MapClaims))
}
