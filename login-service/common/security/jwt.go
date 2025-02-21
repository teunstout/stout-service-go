package security

import (
	"crypto/rsa"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJwt(w http.ResponseWriter, r *http.Request, signKey *rsa.PrivateKey) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.MapClaims{
		"iss": "login-service", // https://golang-jwt.github.io/jwt/usage/signing_methods/
	}

	st, err := t.SignedString(signKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(st))
}
