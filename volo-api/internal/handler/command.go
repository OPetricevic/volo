package handler

import (
	"net/http"
	"strconv"

	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/model"
)

func (h *Handlers) ProcessCommand(w http.ResponseWriter, r *http.Request) {
	var req model.ProcessCommandRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	// Use proper validation (checks empty + max length)
	if errResp := middleware.ValidateTranscript(req.Transcript); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.services.Command.Process(r.Context(), userID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "COMMAND_PARSE_FAILED", "We couldn't process that command. Try again.", err.Error())
		return
	}

	middleware.RecordCommand(result.Action)
	h.respond(w, http.StatusCreated, model.Response{Data: result})
}

func (h *Handlers) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	suggestions, err := h.services.Command.GetSuggestions(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "SUGGESTIONS_FAILED", "Failed to get suggestions.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: suggestions})
}

func (h *Handlers) GetHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	result, err := h.services.Command.GetHistory(r.Context(), userID, page, pageSize)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "HISTORY_FETCH_FAILED", "Failed to fetch history.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) ClearHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	if err := h.services.Command.ClearHistory(r.Context(), userID); err != nil {
		h.respondError(w, http.StatusInternalServerError, "HISTORY_CLEAR_FAILED", "Failed to clear history.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"status": "cleared"}})
}

func (h *Handlers) GetHistoryContext(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	context, err := h.services.Command.GetHistoryContext(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "HISTORY_CONTEXT_FAILED", "Failed to build history context.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: map[string]string{"context": context}})
}

func (h *Handlers) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	settings, err := h.services.Command.GetSettings(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "SETTINGS_FETCH_FAILED", "Failed to get settings.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: settings})
}

func (h *Handlers) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req model.UpdateSettingsRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	settings, err := h.services.Command.UpdateSettings(r.Context(), userID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "SETTINGS_UPDATE_FAILED", "Failed to update settings.", err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: settings})
}
