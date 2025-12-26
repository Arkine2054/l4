package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/arkine/l4/3/internal/calendar"
	httpapi "gitlab.com/arkine/l4/3/internal/http"
	"gitlab.com/arkine/l4/3/internal/logger"
	"gitlab.com/arkine/l4/3/internal/workers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := calendar.NewStorage()
	service := calendar.NewService(storage)
	handler := &httpapi.Handler{Service: service}
	log := logger.NewLogger()

	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: httpapi.LoggingMiddleware(log, router),
	}

	workers.StartReminderPool(ctx, 2, service.ReminderChannel(), service, func(s string) {
		fmt.Println("!!!", s)
	})
	go workers.Cleaner(ctx, service, 10*time.Second)

	fmt.Println("Server started on port", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("HTTP error:", err)
		}
	}()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Shutting down...")
	cancel()
	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(shutdownCtx)
	fmt.Println("Server stopped gracefully")
}
