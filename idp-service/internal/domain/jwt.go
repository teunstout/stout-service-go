package domain

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AppClaims struct {
    // Permissions tell what a user is able to do e.g. ai-service:read or ai-service:read:admin
    Permissions []string `json:"permissions"`
    // The roles at this point is more so a tag,
    // in the permissions table your role relates to permissions.application
    Roles       []string `json:"roles"`
    // The tenant_id is supposed to represent an organization or application. For example,
    // lets say I am owning teun-stout.com and every app including that. Then my tenant would be t_teunstout
    // Now, I start a new project for some internal application called "stout-insurance.com".
    // We now add the tenant_id t_stoutinsurance. The reason is to make sure permissons and roles don't leak
    // to another organization/application.
    TenantID    string   `json:"tenant_id"`
    jwt.RegisteredClaims
}

func CreateJwt(mid int32, signKey *rsa.PrivateKey) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.MapClaims{
		"iss": "idp-server", // https://golang-jwt.github.io/jwt/usage/signing_methods/
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		"jti": uuid.New(),
		"sub": mid,
	}

	return t.SignedString(signKey)
}
