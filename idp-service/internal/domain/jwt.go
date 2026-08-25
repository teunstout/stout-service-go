package domain

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppClaims struct {
	Sub int32 `json:"sub"`

	Permissions []string `json:"permissions"`

	Roles []string `json:"roles"`

	TenantID string `json:"tenant_id"`
	jwt.RegisteredClaims
}

func CreateJwt(mid int32, signKey *rsa.PrivateKey) (string, error) {
	claims := AppClaims{
		Sub:         mid,
		Permissions: []string{},
		Roles:       []string{},
		TenantID:    "",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "idp-server",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			ID:        uuid.New().String(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return t.SignedString(signKey)
}
