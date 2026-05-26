package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/timurturovets/test-backend-rwb/internal/engine"
	"github.com/timurturovets/test-backend-rwb/internal/metrics"
)

type Handler struct {
	engine *engine.Engine
}

func NewHandler(e *engine.Engine) *Handler {
	return &Handler{
		engine: e,
	}
}

type topResponse struct {
	Items []engine.Entry `json:"items"`
	Total int            `json:"total"`
}

func (h *Handler) GetTop(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		metrics.RequestDuration.Observe(time.Since(start).Seconds())
	}()
	metrics.TopRequests.Inc()

	n, err := strconv.Atoi(r.URL.Query().Get("n"))
	if err != nil || n <= 0 {
		n = 10
	}
	if n > 100 {
		n = 100
	}

	items := h.engine.Top(n)
	writeJSON(w, http.StatusOK, topResponse{
		Items: items,
		Total: len(items),
	})
}

func (h *Handler) AddStopword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Word string `json:"word"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Word == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "word is required",
		})
		return
	}

	h.engine.AddToStoplist(body.Word)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveStopword(w http.ResponseWriter, r *http.Request) {
	word := r.PathValue("word")
	if word == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "word is required",
		})
		return
	}

	h.engine.RemoveFromStoplist(word)
	w.WriteHeader(http.StatusNoContent)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
