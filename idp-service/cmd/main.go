package main

import (
	"crypto/rsa"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"stout.dev/idp/internal/handler"
	"stout.dev/idp/internal/repository"
	"stout.dev/idp/internal/usecase"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func main() {
	env := os.Getenv("ENV")
	logger, err := zap.NewDevelopment()
	if env != "" {
		logger, err = zap.NewProduction()
	}
	fatal(err)
	defer logger.Sync()

	connString := os.Getenv("CONNECTION_STRING")
	if connString == "" {
		log.Fatal("CONNECTION_STRING environment variable not set")
	}

	privKeyPath := "/app/keys/app.rsa"
	if _, err := os.Stat(privKeyPath); os.IsNotExist(err) {
		privKeyPath = "../app.rsa"
		logger.Info("Private key not found, using default path: " + privKeyPath)
	}

	signBytes, err := os.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	pubKeyPath := "/app/keys/app.rsa.pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		pubKeyPath = "../app.rsa.pub"
		logger.Info("Public key not found, using default path: " + pubKeyPath)
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
	jwksUsecase := usecase.JwksUsecase(logger, &signKey.PublicKey)

	authHandler := handler.NewAuthenticateHandler(authUsecase, logger)
	registerHandler := handler.NewRegisterHandler(registerUsecase, logger)
	loginHandler := handler.NewLoginHandler(loginUsecase, logger)
	logoutHandler := handler.NewLogoutHandler(logoutUsecase, logger)
	restrictedHandler := handler.NewRestrictedHandler(authUsecase, logger, verifyKey)
	jwksHandler := handler.JwksHandler(jwksUsecase, logger)

	http.HandleFunc("/v1/register", registerHandler.HandleRegister)
	http.HandleFunc("/v1/login", loginHandler.HandleLogin)
	http.HandleFunc("/v1/logout", logoutHandler.HandleLogout)
	http.HandleFunc("/.well-known/jwks.json", jwksHandler.HandleJwksKeys)

	// Handle the session and csrf token
	http.HandleFunc("/v1/auth", authHandler.HandleAuthenticate)

	// This is an example and your spring service has to implement this part.
	http.HandleFunc("/v1/endpoints/restricted/example", restrictedHandler.RestrictedHandler)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
