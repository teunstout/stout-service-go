package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var (
	key []byte
	t   *jwt.Token
	s   string
)

func main() {

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5
	http.HandleFunc("/jwt", func(w http.ResponseWriter, r *http.Request) {
		key = []byte("secret")
		t = jwt.NewWithClaims(jwt.SigningMethodHS256,
			jwt.MapClaims{
				"iss": "jwt-service",
			}) // https://golang-jwt.github.io/jwt/usage/signing_methods/
		s, err := t.SignedString(key)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte(s))
	})

	http.ListenAndServe(":9096", nil)
}
