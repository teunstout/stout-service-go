package handler

import (
	"net/http"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type LogoutHandlerInterface struct {
	logger  *zap.Logger
	usecase *usecase.LogoutUsecaseInterface
}

func NewLogoutHandler(usecase *usecase.LogoutUsecaseInterface, logger *zap.Logger) *LogoutHandlerInterface {
	return &LogoutHandlerInterface{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *LogoutHandlerInterface) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Debug("Method not allowed")
		http.Error(w, domain.MethodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		h.logger.Info("Someone tried logging out without being logged in")
		http.Error(w, domain.ForbiddenMessage, http.StatusForbidden)
		return
	}

	csrfToken := r.Header.Get("x-csrf-token")
	if csrfToken == "" {
		h.logger.Warn("Session token without x-csrf-token")
		http.Error(w, domain.ForbiddenMessage, http.StatusForbidden)
		return
	}

	logoutAllDevices := r.URL.Query().Get("all_devices") == "true"

	h.usecase.Logout(sessionToken.Value, csrfToken, logoutAllDevices)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion logged out successfully"))
}
