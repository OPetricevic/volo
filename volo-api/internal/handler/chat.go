package handler

import (
	"net/http"

	"github.com/volo/volo-api/internal/middleware"
	"github.com/volo/volo-api/internal/model"
	"github.com/volo/volo-api/internal/service"
)

func (h *Handlers) Chat(w http.ResponseWriter, r *http.Request) {
	var req service.ChatRequest
	if err := decode(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body.", err.Error())
		return
	}

	// Use proper validation (checks empty + max length)
	if errResp := middleware.ValidateChatMessage(req.Message); errResp != nil {
		h.respond(w, http.StatusBadRequest, model.Response{Error: errResp})
		return
	}

	userID := middleware.GetUserID(r.Context())
	result, err := h.services.Chat.Process(r.Context(), userID, req)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "CHAT_FAILED",
			"Volo AI couldn't process that. Make sure Ollama is running.",
			err.Error())
		return
	}

	h.respond(w, http.StatusOK, model.Response{Data: result})
}

func (h *Handlers) ChatStatus(w http.ResponseWriter, r *http.Request) {
	running := h.services.Chat.CheckOllamaStatus()
	h.respond(w, http.StatusOK, model.Response{
		Data: map[string]any{
			"ollama_running": running,
			"model":          "gemma2:2b",
		},
	})
}
