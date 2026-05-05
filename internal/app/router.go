package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"pupupu/internal/database"
	"pupupu/internal/logger"
)

func NewRouter(db *database.PSQL) http.Handler {
	mux := http.NewServeMux()

	handleCreate := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var sub database.Subscription

			err := json.NewDecoder(r.Body).Decode(&sub)
			if err != nil {
				logger.Log.Error("Ошибка декодирования", "err", err)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			id, err := db.CreateSub(sub)
			if err != nil {
				logger.Log.Error("Ошибка базы", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"id": %d}\n`, id)
		}
	}

	handleGetOne := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			sub, err := db.GetSubByID(id)

			if err != nil {
				logger.Log.Error("Ошибка получения", "id", id, "err", err)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sub)
		}
	}

	handleUpdate := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			var sub database.Subscription
			if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
				http.Error(w, "Bad JSON", http.StatusBadRequest)
				return
			}

			if err := db.UpdateSub(id, sub); err != nil {
				logger.Log.Error("Update error", "err", err)
				http.Error(w, "Not Found or Error", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status": "updated"}\n`)
		}
	}

	handleDelete := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			if err := db.DeleteSub(id); err != nil {
				logger.Log.Error("Update error", "err", err)
				http.Error(w, "Not Found or Error", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"status": "delete"}\n`)
		}
	}

	handleList := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sub, err := db.GetAllSubs()

			if err != nil {
				logger.Log.Error("Ошибка получения", "err", err)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sub)
		}
	}

	handleTotal := func(db *database.PSQL) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID := r.URL.Query().Get("user_id")
			serviceName := r.URL.Query().Get("service_name")
			period := r.URL.Query().Get("period")

			if userID == "" || serviceName == "" || period == "" {
				http.Error(w, "Missing query parameters (user_id, service_name, period)", http.StatusBadRequest)
				return
			}

			total, err := db.GetTotal(userID, serviceName, period)
			if err != nil {
				logger.Log.Error("Sum error", "err", err)
				http.Error(w, "Error calculating total", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"total_price": %d}\n`, total)
		}
	}

	mux.HandleFunc("POST /subscriptions", handleCreate(db))
	mux.HandleFunc("GET /subscriptions/{id}", handleGetOne(db))
	mux.HandleFunc("PUT /subscriptions/{id}", handleUpdate(db))
	mux.HandleFunc("DELETE /subscriptions/{id}", handleDelete(db))
	mux.HandleFunc("GET /subscriptions", handleList(db))

	mux.HandleFunc("GET /subscriptions/total", handleTotal(db))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	return mux
}
