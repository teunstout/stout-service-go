package domain

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateJwt(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test RSA key: %v", err)
	}

	const accountID int32 = 42
	tokenString, err := CreateJwt(accountID, key)
	if err != nil {
		t.Fatalf("CreateJwt returned error: %v", err)
	}

	claims := &AppClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			t.Fatalf("unexpected signing method: %v", token.Method)
		}
		return &key.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse/verify token: %v", err)
	}
	if !token.Valid {
		t.Fatal("token was parsed but is not valid")
	}

	if claims.Sub != accountID {
		t.Errorf("claims.Sub = %d, want %d", claims.Sub, accountID)
	}
	if claims.Issuer != "idp-server" {
		t.Errorf("claims.Issuer = %q, want %q", claims.Issuer, "idp-server")
	}
	if claims.ID == "" {
		t.Error("claims.ID (jti) is empty, want a uuid")
	}
	if claims.Roles == nil || len(claims.Roles) != 0 {
		t.Errorf("claims.Roles = %v, want empty slice", claims.Roles)
	}
	if claims.Permissions == nil || len(claims.Permissions) != 0 {
		t.Errorf("claims.Permissions = %v, want empty slice", claims.Permissions)
	}
	if claims.TenantID != "" {
		t.Errorf("claims.TenantID = %q, want empty string", claims.TenantID)
	}

	wantExpiry := time.Now().Add(12 * time.Hour)
	gotExpiry := claims.ExpiresAt.Time
	if diff := gotExpiry.Sub(wantExpiry); diff > time.Minute || diff < -time.Minute {
		t.Errorf("claims.ExpiresAt = %v, want ~%v", gotExpiry, wantExpiry)
	}
}

func TestCreateJwtSubIsJSONNumber(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate test RSA key: %v", err)
	}

	tokenString, err := CreateJwt(7, key)
	if err != nil {
		t.Fatalf("CreateJwt returned error: %v", err)
	}

	// jisho-service's AuthMiddleware parses claims generically as
	// jwt.MapClaims and requires claims["sub"] to type-assert to float64 -
	// verify that contract still holds after switching CreateJwt to AppClaims.
	mapClaims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, mapClaims, func(token *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})
	if err != nil {
		t.Fatalf("failed to parse token into MapClaims: %v", err)
	}

	sub, ok := mapClaims["sub"].(float64)
	if !ok {
		t.Fatalf("claims[\"sub\"] is not a float64, got %T", mapClaims["sub"])
	}
	if sub != 7 {
		t.Errorf("claims[\"sub\"] = %v, want 7", sub)
	}
}
