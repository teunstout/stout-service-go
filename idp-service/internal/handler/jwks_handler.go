package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type JwksHandlerInterface struct {
	usecase *usecase.JwksUsecaseInterface
	logger  *zap.Logger
}

func JwksHandler(usecase *usecase.JwksUsecaseInterface, logger *zap.Logger) *JwksHandlerInterface {
	return &JwksHandlerInterface{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *JwksHandlerInterface) HandleJwksKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	jwks := h.usecase.JwksKeys()
	h.logger.Debug("Requested JWKS Keys", zap.Any("keys", jwks))

	jsonResponse, err := json.Marshal(jwks)
	if err != nil {
		h.logger.Info("Creating Json response failed", zap.Any("body", jwks))
		http.Error(w, domain.InternalServerErrorMessage, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
