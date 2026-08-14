package handler

import (
	"encoding/json"

	"net/http"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type LoginHandlerInterface struct {
	logger        *zap.Logger
	usecase       *usecase.LoginUsecaseInterface
	secureCookies bool
}

func NewLoginHandler(usecase *usecase.LoginUsecaseInterface, logger *zap.Logger, secureCookies bool) *LoginHandlerInterface {
	return &LoginHandlerInterface{
		usecase:       usecase,
		logger:        logger,
		secureCookies: secureCookies,
	}
}

func (h *LoginHandlerInterface) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Debug("Method not allowed", zap.String("method", r.Method))
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	var loginData domain.LoginBody
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		h.logger.Debug("Invalid JSON payload", zap.Error(err))
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		h.logger.Debug("Missing username or password", zap.String("username", loginData.Username))
		writeJSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	result, err := h.usecase.Login(loginData.Username, loginData.Password)
	if err != nil {
		h.logger.Info("Login failed", zap.String("username", loginData.Username), zap.Error(err))
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    result.SessionToken,
		Path:     "/",
		Expires:  result.ExpiresAt,
		HttpOnly: true,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	})

	body := map[string]string{"jwt": result.Jwt, "csrfToken": result.CsrfToken}
	jsonResponse, err := json.Marshal(body)
	if err != nil {
		h.logger.Info("Creating Json response failed", zap.Any("body", body), zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
