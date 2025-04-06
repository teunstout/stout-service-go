package security

import (
	"crypto/rsa"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJwt(w http.ResponseWriter, r *http.Request, mid int32, signKey *rsa.PrivateKey) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.MapClaims{
		"iss": "login-service", // https://golang-jwt.github.io/jwt/usage/signing_methods/
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"jti": uuid.New(),
		"sub": mid,
	}

	st, err := t.SignedString(signKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the JWT as an HttpOnly cookie with a short name
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    st,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}
