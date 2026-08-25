package domain

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwksKeys struct {
	Keys []Jwks `json:"keys"`
}

type Jwks struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	X5c []string `json:"x5c,omitempty"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t,omitempty"`
}

func CreateJwksKeys(publicKey *rsa.PublicKey) JwksKeys {
	n, e := encodeJWKParams(publicKey)

	key := Jwks{
		Alg: jwt.SigningMethodRS256.Alg(),
		Kty: "RSA",
		Use: "sig",
		N:   n,
		E:   e,
		Kid: fmt.Sprintf("idp-public-key-%d", time.Now().Unix()),
	}
	return JwksKeys{
		Keys: []Jwks{key},
	}
}

func encodeJWKParams(pub *rsa.PublicKey) (n string, e string) {

	n = base64.RawURLEncoding.EncodeToString(pub.N.Bytes())

	eBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(eBytes, uint64(pub.E))
	eBytes = bytes.TrimLeft(eBytes, "\x00")
	e = base64.RawURLEncoding.EncodeToString(eBytes)

	return n, e
}
