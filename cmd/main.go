package main

import (
	"net/http"
	"time"

	"pupupu/internal/app"
	"pupupu/internal/config"
	"pupupu/internal/database"
	"pupupu/internal/logger"
)

func main() {
	cfg := config.New()

	logger.Init("pupupu.log")
	defer logger.Close()

	db := database.Init(cfg)
	defer db.Close()

	router := app.NewRouter(db)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Info("Gateway is running", "addr", 8080)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("Server stopped with error", "err", err)
	}
}
