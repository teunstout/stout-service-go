package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"stout.dev/login/internal/handler"
	"stout.dev/login/internal/repository"
	"stout.dev/login/internal/usecase"
	logruslogger "stout.dev/login/pkg/logger"
)

const (
	DATABASE_URL = "user=golang password=golang host=localhost port=5432 dbname=postgres sslmode=disable"
)

var (
	verifyKey  *rsa.PublicKey
	signKey    *rsa.PrivateKey
	serverPort int
)

func main() {
	logger := logruslogger.NewLogrusLogger()
	connString := os.Getenv("CONNECTION_STRING")
	if connString == "" {
		connString = DATABASE_URL
	}

	privKeyPath := "/app/keys/app.rsa"
	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		fmt.Println("Private key not found, using default path")
		privKeyPath = "../app.rsa"
	}

	signBytes, err := os.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	pubKeyPath := "/app/keys/app.rsa.pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		fmt.Println("Public key not found, using default path")
		pubKeyPath = "../app.rsa.pub"
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	accountRepo := repository.NewAccountRepository(connString, logger)
	sessionRepo := repository.NewSessionTokenRepository(connString, logger)
	csrfRepo := repository.NewCsrfRepository(connString, logger)

	authUsecase := usecase.NewAuthenticationUsecase(sessionRepo, csrfRepo, logger)
	loginUsecase := usecase.NewLoginUsecase(accountRepo, sessionRepo, csrfRepo, logger, signKey)
	registerUsecase := usecase.NewRegisterUsecase(accountRepo, logger, signKey)
	logoutUsecase := usecase.NewLogoutUsecase(sessionRepo, csrfRepo, logger)

	authHandler := handler.NewAuthenticateHandler(authUsecase, logger)
	registerHandler := handler.NewRegisterHandler(registerUsecase, logger)
	loginHandler := handler.NewLoginHandler(loginUsecase, logger)
	logoutHandler := handler.NewLogoutHandler(logoutUsecase, logger)
	restrictedHandler := handler.NewRestrictedHandler(authUsecase, logger, verifyKey)

	http.HandleFunc("/v1/register", registerHandler.HandleRegister)
	http.HandleFunc("/v1/login", loginHandler.HandleLogin)
	http.HandleFunc("/v1/logout", logoutHandler.HandleLogout)
	http.HandleFunc("/v1/authenticate", authHandler.HandleAuthenticate)
	http.HandleFunc("/v1/endpoints/restricted", restrictedHandler.RestrictedHandler)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
