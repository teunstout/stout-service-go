package handler

import (
	"encoding/json"

	"net/http"
	"time"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type LoginHandlerInterface struct {
	logger  *zap.Logger
	usecase *usecase.LoginUsecaseInterface
}

func NewLoginHandler(usecase *usecase.LoginUsecaseInterface, logger *zap.Logger) *LoginHandlerInterface {
	return &LoginHandlerInterface{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *LoginHandlerInterface) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.logger.Debug("Method not allowed", zap.String("method", r.Method))
		http.Error(w, domain.MethodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	var loginData domain.LoginBody
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		h.logger.Debug("Invalid JSON payload", zap.Error(err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if loginData.Username == "" || loginData.Password == "" {
		h.logger.Debug("Invalid login attempt",
			zap.String("username", loginData.Username),
			zap.String("password", loginData.Password),
		)
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	session, csrf, jwt, err := h.usecase.Login(loginData.Username, loginData.Password)
	if err != nil {
		h.logger.Info("Login failed", zap.String("username", loginData.Username))
		http.Error(w, domain.UnauthorizedMessage, http.StatusUnauthorized)
		return
	}

	exprDate := time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     "JWT",
		Value:    jwt,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  exprDate,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrf.Value,
		Path:     "/",
		Expires:  exprDate,
		HttpOnly: false,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Value,
		Path:     "/",
		Expires:  exprDate,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}
