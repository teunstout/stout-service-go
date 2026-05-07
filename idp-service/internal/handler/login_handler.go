package handler

import (
	"encoding/json"

	"net/http"

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

	body := map[string]string{"jwt": jwt, "csrfToken": csrf.Value, "sessionToken": session.Value}
	jsonResponse, err := json.Marshal(body)
	if err != nil {
		h.logger.Info("Creating Json response failed", zap.Any("body", body))
		http.Error(w, domain.InternalServerErrorMessage, http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
