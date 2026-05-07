package usecase

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JwksUsecaseInterface struct {
	logger            *zap.Logger
	publicKey         *rsa.PublicKey
	publicKeyEncodedE string
	publicKeyEncodedN string
}

func JwksUsecase(logger *zap.Logger, publicKey *rsa.PublicKey) *JwksUsecaseInterface {
	n, e := encodeJWKParams(publicKey)
	logger.Info("Public key fields", zap.String("N", n), zap.String("e", e))

	return &JwksUsecaseInterface{
		logger:            logger,
		publicKey:         publicKey,
		publicKeyEncodedE: e,
		publicKeyEncodedN: n,
	}
}

type JwksKeys struct {
	Keys []Jwks `json:"keys"`
}

// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-key-set-properties
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

func (u *JwksUsecaseInterface) JwksKeys() (k JwksKeys) {
	key := Jwks{
		Alg: jwt.SigningMethodRS256.Alg(),
		Kty: "RSA",
		Use: "sig",
		N:   u.publicKeyEncodedN,
		E:   u.publicKeyEncodedE,
		Kid: fmt.Sprintf("idp-public-key-%d", time.Now().Unix()),
	}
	return JwksKeys{
		Keys: []Jwks{key},
	}
}

func encodeJWKParams(pub *rsa.PublicKey) (n string, e string) {
	// N: big.Int has a Bytes() method that returns big-endian bytes
	n = base64.RawURLEncoding.EncodeToString(pub.N.Bytes())

	// E: convert int to big-endian bytes, trim leading zeros
	eBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(eBytes, uint64(pub.E))
	eBytes = bytes.TrimLeft(eBytes, "\x00")
	e = base64.RawURLEncoding.EncodeToString(eBytes)

	return n, e
}
