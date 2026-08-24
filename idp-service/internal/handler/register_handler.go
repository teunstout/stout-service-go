package handler

import (
	"encoding/json"
	"errors"

	"net/http"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/usecase"
)

type RegisterHandlerInterface struct {
	logger  *zap.Logger
	usecase *usecase.RegisterUsecaseInterface
}

func NewRegisterHandler(usecase *usecase.RegisterUsecaseInterface, logger *zap.Logger) *RegisterHandlerInterface {
	return &RegisterHandlerInterface{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *RegisterHandlerInterface) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	var loginData domain.LoginBody

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	username := loginData.Username
	password := loginData.Password

	if username == "" || password == "" {
		writeJSONError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	if err := h.usecase.Register(username, password); err != nil {
		if errors.Is(err, domain.ErrAccountAlreadyExists) {
			writeJSONError(w, http.StatusConflict, "Account already exists")
			return
		}
		h.logger.Error("Register failed", zap.String("username", username), zap.Error(err))
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Minion registered successfully"))
}
