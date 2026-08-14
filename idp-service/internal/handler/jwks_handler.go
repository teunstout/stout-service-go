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
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	jwks := h.usecase.JwksKeys()

	jsonResponse, err := json.Marshal(jwks)
	if err != nil {
		h.logger.Info("Creating Json response failed", zap.Any("body", jwks))
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
