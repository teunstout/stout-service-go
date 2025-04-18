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
		// Retrieve the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Debug("Authorization header missing", nil)
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Ensure the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Debug("Invalid Authorization header format", nil)
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is RS256
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

		// Extract claims and retrieve the "sub" (member ID)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.Warn("Invalid token claims", nil)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		sub, ok := claims["sub"].(float64) // JWT encodes numbers as float64
		if !ok {
			http.Error(w, "Member ID (sub) missing in token", http.StatusUnauthorized)
			return
		}

		// Convert "sub" to int32 and add it to the request context
		mid := int32(sub)
		ctx := context.WithValue(r.Context(), "mid", mid)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}
