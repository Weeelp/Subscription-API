package handler

import (
	"encoding/json"
	"net/http"
	"pupupu/internal/models"
	"pupupu/internal/service"
	"strconv"
)

type SubscriptionHandler struct {
	srv *service.SubService
}

func NewSubscriptionHandler(srv *service.SubService) *SubscriptionHandler {
	return &SubscriptionHandler{srv: srv}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	id, err := h.srv.CreateNewSubscription(r.Context(), sub)
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}

	h.sendJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *SubscriptionHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	sub, err := h.srv.GetSubByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	h.sendJSON(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	if err := h.srv.UpdateSub(r.Context(), id, sub); err != nil {
		http.Error(w, "Not Found or Error", http.StatusNotFound)
		return
	}

	h.sendJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.srv.DeleteSub(r.Context(), id); err != nil {
		http.Error(w, "Not Found or Error", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	subs, err := h.srv.GetAllSubs(r.Context())
	if err != nil {
		http.Error(w, "Error fetching list", http.StatusInternalServerError)
		return
	}

	h.sendJSON(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) Total(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	period := r.URL.Query().Get("period")

	total, err := h.srv.GetTotal(r.Context(), userID, serviceName, period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.sendJSON(w, http.StatusOK, map[string]int{"total_price": total})
}

func (h *SubscriptionHandler) parseID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func (h *SubscriptionHandler) sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
