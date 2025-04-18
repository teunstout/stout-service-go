package handler

import (
	"encoding/json"
	"net/http"

	"stout.dev/jisho/internal/domain"
	"stout.dev/jisho/internal/usecase"
)

type JishoHandlerInterface struct {
	usecase *usecase.JishoUsecaseInterface
	logger  domain.Logger
}

func NewJishoHandler(u *usecase.JishoUsecaseInterface, d domain.Logger) *JishoHandlerInterface {
	return &JishoHandlerInterface{usecase: u, logger: d}
}

func (h *JishoHandlerInterface) SearchJisho(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.logger.Debug("Invalid method", map[string]interface{}{"method": r.Method})
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	kw := r.URL.Query().Get("keyword")
	if kw == "" {
		h.logger.Debug("Keyword is empty", map[string]interface{}{"keyword": kw})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, _ := h.usecase.SearchJisho(kw, r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
