package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"stout.dev/content/internal/domain"
	"stout.dev/content/internal/middleware"
	"stout.dev/content/internal/usecase"
)

type TranslationHandler struct {
	usecase *usecase.TranslationUsecase
}

func NewTranslationHandler(u *usecase.TranslationUsecase) *TranslationHandler {
	return &TranslationHandler{usecase: u}
}

func (h *TranslationHandler) HandleSyncList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	var req domain.SyncListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		writeJSONError(w, http.StatusBadRequest, "name is required")
		return
	}

	result, err := h.usecase.SyncList(r.Context(), accountID, req)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *TranslationHandler) HandleGetLists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	result, err := h.usecase.GetLists(r.Context(), accountID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *TranslationHandler) HandleDeleteList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	var req domain.DeleteListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	deleted, err := h.usecase.DeleteList(r.Context(), accountID, req.ID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}
	if !deleted {
		writeJSONError(w, http.StatusNotFound, "List not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.DeleteListResult{ID: req.ID})
}

func (h *TranslationHandler) HandleDeleteEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, domain.MethodNotAllowedMessage)
		return
	}

	accountID, ok := middleware.AccountIDFromContext(r.Context())
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, domain.UnauthorizedMessage)
		return
	}

	var req domain.DeleteEntriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	deletedIDs, err := h.usecase.DeleteEntries(r.Context(), accountID, req.IDs)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, domain.InternalServerErrorMessage)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.DeleteEntriesResult{DeletedIDs: deletedIDs})
}
