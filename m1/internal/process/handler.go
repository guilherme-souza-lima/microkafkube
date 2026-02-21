package process

import (
	"encoding/json"
	"io"
	"net/http"
)

type HandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || r.ContentLength == 0 {
		h.response(w, http.StatusBadRequest, map[string]string{"error": "request body is empty"})
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.response(w, http.StatusBadRequest, map[string]string{"error": "failed to read body"})
		return
	}

	response, err := h.service.Register(r.Context(), body)
	if err != nil {
		h.response(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	h.response(w, http.StatusCreated, response)
}

func (h *Handler) response(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
