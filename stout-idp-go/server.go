package main

import (
	"context"
	"encoding/json"
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
)

const DATABASE_URL = "user=golang password=golang host=127.0.0.1 port=5432 dbname=golang sslmode=disable"

func main() {
	// Connect to the database and create a new adapter
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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

	manager := manage.NewDefaultManager()
	manager.MapTokenStorage(tokenStore)
	manager.MapClientStorage(clientStore)

	// Create an instance of the server
	srv := server.NewDefaultServer(manager)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)
	})

	http.HandleFunc("/v1/clients", func(w http.ResponseWriter, r *http.Request) {
		handleClientRequests(w, r, conn)
	})

	http.HandleFunc("/v1/tokens", func(w http.ResponseWriter, r *http.Request) {
		handleClientTokens(w, r, conn)
	})

	log.Println("Started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":9096", nil))
}

func handleClientTokens(w http.ResponseWriter, _ *http.Request, conn *pgx.Conn) {
	rows, err := conn.Query(context.Background(), "SELECT code, access, refresh FROM oauth2_tokens ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var tokens []models.Token
	for rows.Next() {
		var token models.Token
		err := rows.Scan(&token.Code, &token.Access, &token.Refresh)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error scanning client", http.StatusInternalServerError)
			return
		}
		tokens = append(tokens, token)
	}

	// Return the clients as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)
}

func handleClientRequests(w http.ResponseWriter, _ *http.Request, conn *pgx.Conn) {
	rows, err := conn.Query(context.Background(), "SELECT id, secret, domain FROM oauth2_clients ORDER BY id ASC")
	if err != nil {
		http.Error(w, "Invalid client_id", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		err := rows.Scan(&client.ID, &client.Secret, &client.Domain)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error scanning client", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	// Return the clients as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)
}
