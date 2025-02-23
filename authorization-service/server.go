package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/jackc/pgx/v4"
	pg "github.com/vgarvardt/go-oauth2-pg/v4"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
	"stout.dev/authorization/database"
	"stout.dev/authorization/login"
)

func main() {
	conn, err := database.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	if err = database.CreateTables(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create tables: %v\n", err)
		os.Exit(1)
	}

	adapter := pgx4adapter.NewConn(conn)

	tokenStore, err := pg.NewTokenStore(adapter, pg.WithTokenStoreGCInterval(time.Minute))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create token store: %v\n", err)
		os.Exit(1)
	}
	defer tokenStore.Close()

	clientStore, err := pg.NewClientStore(adapter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create client store: %v\n", err)
		os.Exit(1)
	}

	clientStore.Create(&models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "localhost",
	})

	// Manager creates access tokens and refresh tokens using UUID and MD5 as base64
	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// Set some server configuration
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	srv.SetUserAuthorizationHandler(func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
		return userAuthorizeHandler(conn, w, r)
	})

	// Set the internal and response error handlers
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/oauth/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	log.Println("Started on http://localhost:9096")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func userAuthorizeHandler(conn *pgx.Conn, w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	userID, err = login.AuthenticateUser(conn, username, password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	return userID, nil
}
