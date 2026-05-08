package handler

import (
	"encoding/json"

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
		http.Error(w, domain.MethodNotAllowedMessage, http.StatusMethodNotAllowed)
		return
	}

	var loginData domain.LoginBody

	// Decode the JSON body
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	username := loginData.Username
	password := loginData.Password

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if err := h.usecase.Register(username, password); err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Minion registered successfully"))
}
