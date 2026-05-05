package handler

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"go.uber.org/zap"
	"stout.dev/idp/internal/usecase"
)

type RestrictedHandlerInterface struct {
	logger    *zap.Logger
	usecase   *usecase.AuthenticationUsecaseInterface
	verifyKey *rsa.PublicKey
}

func NewRestrictedHandler(usecase *usecase.AuthenticationUsecaseInterface, logger *zap.Logger, verifyKey *rsa.PublicKey) *RestrictedHandlerInterface {
	return &RestrictedHandlerInterface{
		usecase:   usecase,
		logger:    logger,
		verifyKey: verifyKey,
	}
}

func (h *RestrictedHandlerInterface) RestrictedHandler(w http.ResponseWriter, r *http.Request) {
	// Get token from request
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return h.verifyKey, nil
	}, request.WithClaims(&jwt.MapClaims{}))

	// If the token is missing or invalid, return error
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token:", err)
		return
	}

	// Token is valid
	fmt.Fprintln(w, "Welcome,", token.Claims.(*jwt.MapClaims))
}
