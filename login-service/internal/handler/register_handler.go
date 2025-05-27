package handler

import (
	"encoding/json"

	"net/http"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/usecase"
)

type RegisterHandlerInterface struct {
	logger  domain.Logger
	usecase *usecase.RegisterUsecaseInterface
}

func NewRegisterHandler(usecase *usecase.RegisterUsecaseInterface, logger domain.Logger) *RegisterHandlerInterface {
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Minion registered successfully"))
}
