package app

import (
	"net/http"

	"pupupu/internal/handler"
)

func NewRouter(h *handler.SubscriptionHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /subscriptions", h.Create)
	mux.HandleFunc("GET /subscriptions/{id}", h.GetOne)
	mux.HandleFunc("PUT /subscriptions/{id}", h.Update)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.Delete)
	mux.HandleFunc("GET /subscriptions", h.List)

	mux.HandleFunc("GET /subscriptions/total", h.Total)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	return mux
}
