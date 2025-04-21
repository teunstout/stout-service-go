package handler

import (
	"encoding/json"
	"example-service/internal/domain"
	"example-service/internal/usecase"
	"net/http"
)

type ExampleHandler struct {
	usecase *usecase.ExampleUsecase
}

func NewExampleHandler(u *usecase.ExampleUsecase) *ExampleHandler {
	return &ExampleHandler{usecase: u}
}

func (h *ExampleHandler) HandleExample(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input domain.ExampleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.usecase.ProcessExample(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
