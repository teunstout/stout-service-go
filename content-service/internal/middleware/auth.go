package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"stout.dev/content/internal/domain"
)

type contextKey string

const accountIDContextKey contextKey = "accountID"

func AuthMiddleware(next http.HandlerFunc, verifyKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, http.ErrAbortHandler
			}
			return verifyKey, nil
		})
		if err != nil || !token.Valid {
			writeUnauthorized(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeUnauthorized(w)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			writeUnauthorized(w)
			return
		}

		accountID := int32(sub)
		ctx := context.WithValue(r.Context(), accountIDContextKey, accountID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AccountIDFromContext(ctx context.Context) (int32, bool) {
	accountID, ok := ctx.Value(accountIDContextKey).(int32)
	return accountID, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": domain.UnauthorizedMessage})
}
