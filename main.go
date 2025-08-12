// cmd/api/main.go
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	addr := getenv("PORT", "8000")

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Recoverer)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Example resource
	r.Route("/slack", func(r chi.Router) {
		r.Get("/event-calendar", getEventCalendarLink)
		r.Post("/event-calendar", getEventCalendarLinkForSlack)
	})

	srv := &http.Server{Addr: ":" + addr, Handler: r}
	go func() {
		logger.Info("server starting", "addr", srv.Addr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("server shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
	logger.Info("bye")
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getEventCalendarLink(w http.ResponseWriter, r *http.Request) {
	// Simulate fetching a calendar link
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"calendarLink":"https://secretive-jade-5a0.notion.site/Event-Calendar-Q3-Q4-2025-24c164277afd80418ab9e78af151fded"}`))
}

func getEventCalendarLinkForSlack(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The event calendar link is: https://secretive-jade-5a0.notion.site/Event-Calendar-Q3-Q4-2025-24c164277afd80418ab9e78af151fded"))
}
