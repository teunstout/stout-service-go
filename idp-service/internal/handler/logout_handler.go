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
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		h.logger.Info("Someone tried logging out without being logged in")
		writeJSONError(w, http.StatusForbidden, domain.ForbiddenMessage)
		return
	}

	csrfToken := r.Header.Get("x-csrf-token")
	if csrfToken == "" {
		h.logger.Warn("Session token without x-csrf-token")
		writeJSONError(w, http.StatusForbidden, domain.ForbiddenMessage)
		return
	}

	logoutAllDevices := r.URL.Query().Get("all_devices") == "true"

	h.usecase.Logout(sessionToken.Value, csrfToken, logoutAllDevices)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Minion logged out successfully"))
}
