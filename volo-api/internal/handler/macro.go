package handler

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/repository"
)

type macroRequest struct {
	TriggerPhrase string                   `json:"trigger_phrase"`
	Name          string                   `json:"name"`
	Actions       []repository.MacroAction `json:"actions"`
	Enabled       *bool                    `json:"enabled,omitempty"`
}

// ListMacros returns all macros for the authenticated user.
func (h *Handlers) ListMacros(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	macros, err := h.services.Repos().Macro.List(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "MACRO_LIST_FAILED", "Failed to load macros.", err.Error())
		return
	}
	if macros == nil {
		macros = []repository.Macro{}
	}
	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]interface{}{
			"macros": macros,
			"limit":  3,
			"count":  len(macros),
		},
	})
}

// ListEnabledMacros returns only enabled macros (for extension sync).
func (h *Handlers) ListEnabledMacros(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	macros, err := h.services.Repos().Macro.ListEnabled(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "MACRO_LIST_FAILED", "Failed to load macros.", err.Error())
		return
	}
	if macros == nil {
		macros = []repository.Macro{}
	}
	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]interface{}{"macros": macros},
	})
}

// CreateMacro adds a new macro.
func (h *Handlers) CreateMacro(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req macroRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_BODY", "Invalid request body.", "")
		return
	}

	req.TriggerPhrase = strings.TrimSpace(strings.ToLower(req.TriggerPhrase))
	req.Name = strings.TrimSpace(req.Name)

	if req.TriggerPhrase == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_TRIGGER", "Trigger phrase is required.", "")
		return
	}
	if req.Name == "" {
		h.respondError(w, http.StatusBadRequest, "MISSING_NAME", "Macro name is required.", "")
		return
	}
	if len(req.Actions) == 0 {
		h.respondError(w, http.StatusBadRequest, "MISSING_ACTIONS", "At least one action is required.", "")
		return
	}
	if len(req.Actions) > 3 {
		h.respondError(w, http.StatusBadRequest, "TOO_MANY_ACTIONS", "Maximum 3 actions per macro.", "")
		return
	}

	macro, err := h.services.Repos().Macro.Create(r.Context(), userID, req.TriggerPhrase, req.Name, req.Actions)
	if err != nil {
		if strings.Contains(err.Error(), "macro limit reached") {
			h.respondError(w, http.StatusForbidden, "MACRO_LIMIT", "Maximum of 10 macros allowed.", "")
			return
		}
		h.respondError(w, http.StatusInternalServerError, "MACRO_CREATE_FAILED", "Failed to create macro.", err.Error())
		return
	}

	h.respond(w, http.StatusCreated, model.Response{Data: macro})
}

// UpdateMacro modifies an existing macro.
func (h *Handlers) UpdateMacro(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	macroID := chi.URLParam(r, "id")

	var req macroRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_BODY", "Invalid request body.", "")
		return
	}

	req.TriggerPhrase = strings.TrimSpace(strings.ToLower(req.TriggerPhrase))
	req.Name = strings.TrimSpace(req.Name)

	if req.TriggerPhrase == "" || req.Name == "" || len(req.Actions) == 0 {
		h.respondError(w, http.StatusBadRequest, "MISSING_FIELDS", "Trigger phrase, name, and actions are required.", "")
		return
	}

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	if err := h.services.Repos().Macro.Update(r.Context(), userID, macroID, req.TriggerPhrase, req.Name, req.Actions, enabled); err != nil {
		h.respondError(w, http.StatusInternalServerError, "MACRO_UPDATE_FAILED", "Failed to update macro.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "updated"}})
}

// DeleteMacro soft-deletes a macro.
func (h *Handlers) DeleteMacro(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	macroID := chi.URLParam(r, "id")

	if err := h.services.Repos().Macro.Delete(r.Context(), userID, macroID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "MACRO_DELETE_FAILED", "Failed to delete macro.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "deleted"}})
}
