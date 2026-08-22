package app

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"stout.dev/content/internal/handler"
	"stout.dev/content/internal/middleware"
	"stout.dev/content/internal/repository"
	"stout.dev/content/internal/usecase"
)

func NewApp() {
	connString := os.Getenv("CONNECTION_STRING")
	if connString == "" {
		log.Fatal("CONNECTION_STRING environment variable not set")
	}

	pool, err := pgxpool.New(context.Background(), connString)
	fatal(err)
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	verifyKey := loadVerifyKey()

	translationRepo := repository.NewTranslationRepository(pool)
	translationUsecase := usecase.NewTranslationUsecase(translationRepo)
	translationHandler := handler.NewTranslationHandler(translationUsecase)

	http.HandleFunc(
		"/v1/content/translation-lists/sync",
		middleware.AuthMiddleware(translationHandler.HandleSyncList, verifyKey),
	)
	http.HandleFunc(
		"/v1/content/translation-lists",
		middleware.AuthMiddleware(translationHandler.HandleGetLists, verifyKey),
	)
	http.HandleFunc(
		"/v1/content/translation-lists/delete",
		middleware.AuthMiddleware(translationHandler.HandleDeleteList, verifyKey),
	)

	log.Println("Started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadVerifyKey() *rsa.PublicKey {
	pubKeyPath := "/app/keys/app.rsa.pub"
	if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
		pubKeyPath = "../app.rsa.pub"
		log.Println("Public key not found, using default path: " + pubKeyPath)
	}

	verifyBytes, err := os.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	return verifyKey
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
