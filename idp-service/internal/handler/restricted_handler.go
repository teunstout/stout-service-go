package handler

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
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
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {

		return h.verifyKey, nil
	}, request.WithClaims(&jwt.MapClaims{}))

	if err != nil {
		h.logger.Info("Invalid token", zap.Error(err))
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	h.logger.Info("Token", zap.Any("Claims", token.Claims.(*jwt.MapClaims)))
	jsonRes, err := json.Marshal(token.Claims.(*jwt.MapClaims))
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
