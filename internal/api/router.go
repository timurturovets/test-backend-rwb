package api

import "net/http"

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/top", h.GetTop)
	mux.HandleFunc("POST /api/v1/stoplist", h.AddStopword)
	mux.HandleFunc("DELETE /api/v1/stoplist/{word}", h.RemoveStopword)

	return mux
}
