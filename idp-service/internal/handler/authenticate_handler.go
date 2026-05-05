package handler

import (
	"net/http"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type AuthenticateHandlerInterface struct {
	logger  *zap.Logger
	usecase *usecase.AuthenticationUsecaseInterface
}

func NewAuthenticateHandler(usecase *usecase.AuthenticationUsecaseInterface, logger *zap.Logger) *AuthenticateHandlerInterface {
	return &AuthenticateHandlerInterface{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *AuthenticateHandlerInterface) HandleAuthenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, domain.MethodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, domain.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	csrfToken := r.Header.Get("x-csrf-token")
	if csrfToken == "" {
		http.Error(w, domain.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	if err := h.usecase.Register(sessionToken.Value, csrfToken); err != nil {
		h.logger.Error("Authentication failed",
			zap.String("sessionToken", sessionToken.Value),
			zap.String("csrfToken", csrfToken),
			zap.Error(err),
		)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion is logged in"))
}
