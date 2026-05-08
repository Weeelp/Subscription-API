package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pupupu/internal/app"
	"pupupu/internal/config"
	"pupupu/internal/handler"
	"pupupu/internal/logger"
	"pupupu/internal/repository"
	"pupupu/internal/service"
)

func main() {
	cfg := config.New()

	logger.Init("pupupu.log")
	defer logger.Close()

	pqRepo := repository.Init(cfg)
	defer pqRepo.Close()

	subService := service.NewSubService(pqRepo, logger.Log)
	subHandler := handler.NewSubscriptionHandler(subService)
	router := app.NewRouter(subHandler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Log.Info("Gateway is running", "addr", 8080)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server error", "err", err)
		}
	}()

	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", "err", err)
	}

	logger.Log.Info("Server exited properly")
}
