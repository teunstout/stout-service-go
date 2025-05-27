package handler

import (
	"net/http"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/usecase"
)

type AuthenticateHandlerInterface struct {
	logger  domain.Logger
	usecase *usecase.AuthenticationUsecaseInterface
}

func NewAuthenticateHandler(usecase *usecase.AuthenticationUsecaseInterface, logger domain.Logger) *AuthenticateHandlerInterface {
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
		h.logger.Error("Authentication failed", map[string]interface{}{
			"sessionToken": sessionToken.Value,
			"csrfToken":    csrfToken,
			"error":        err,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion is logged in"))
}
