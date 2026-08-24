package main

import (
	"crypto/rsa"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	logruslogger "stout.dev/jisho/internal/adapters/logger"
	jishoclient "stout.dev/jisho/internal/external/jishoClient"
	"stout.dev/jisho/internal/handler"
	"stout.dev/jisho/internal/middleware"
	"stout.dev/jisho/internal/repository"
	"stout.dev/jisho/internal/usecase"
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

	pubKeyPath := "/app/keys/app.rsa.pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		logger.Warn("Public key not found, using default path", nil)
		pubKeyPath = `../app.rsa.pub`
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	mux := http.NewServeMux()

	jishoRepository := repository.NewJishoRepository(connString, logger)
	jishoClient := jishoclient.NewJishoUsecase(logger)
	jishoUsecase := usecase.NewJishoUsecase(jishoClient, jishoRepository, logger)
	jishoHandler := handler.NewJishoHandler(jishoUsecase, logger)

	mux.HandleFunc("/v1/search", middleware.AuthMiddleware(jishoHandler.SearchJisho, logger, verifyKey))
	http.ListenAndServe(":8080", mux)
	logger.Info("Server ready and listening!", nil)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
