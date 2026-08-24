package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"stout.dev/jisho/internal/domain"
)

func AuthMiddleware(next http.HandlerFunc, logger domain.Logger, verifyKey *rsa.PublicKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Debug("Authorization header missing", nil)
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Debug("Invalid Authorization header format", nil)
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
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
			logger.Warn("Invalid or expired token", map[string]interface{}{"error": err})
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("Invalid token claims", nil)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Member ID (sub) missing in token", http.StatusUnauthorized)
			return
		}

		mid := int32(sub)
		ctx := context.WithValue(r.Context(), "mid", mid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
