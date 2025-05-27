package domain

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func CreateJwt(mid int32, signKey *rsa.PrivateKey) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.MapClaims{
		"iss": "login-service", // https://golang-jwt.github.io/jwt/usage/signing_methods/
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"jti": uuid.New(),
		"sub": mid,
	}

	return t.SignedString(signKey)
}
