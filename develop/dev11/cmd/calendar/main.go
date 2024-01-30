package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Karanth1r3/l_2/develop/dev11/httpapi"
	"github.com/Karanth1r3/l_2/develop/dev11/internal/config"
	"github.com/Karanth1r3/l_2/develop/dev11/internal/service"
	"github.com/Karanth1r3/l_2/develop/dev11/internal/storage"
)

func configureRoute(h *httpapi.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/create_event", http.HandlerFunc(h.CreateEvent))
	mux.Handle("/update_event", http.HandlerFunc(h.UpdateEvent))
	mux.Handle("/delete_event", http.HandlerFunc(h.DeleteEvent))
	mux.Handle("/events_for_day", http.HandlerFunc(h.GetEventsForDay))
	mux.Handle("/events_for_week", http.HandlerFunc(h.GetEventsForWeek))
	mux.Handle("/events_for_month", http.HandlerFunc(h.GetEventsForMonth))
	return mux
}

func main() {

	cfg, err := config.Read("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Initializing storage, service containing storage & handler object which contains all handler methods
	strg := storage.NewInMemStorage()
	srvc := service.New(strg)
	handler := httpapi.NewHandler(srvc)

	// Configuring router
	mux := configureRoute(handler)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	server := http.Server{
		Addr: fmt.Sprintf(":%d", cfg.HTTP.Port),
		// Middleware logger is set as first handler for all ... individual handlers in router this way.
		// Because logger should track all requests. If Middleware is required for individual handler - it can be passed to mux.Handle wrapping original handler
		Handler: httpapi.NewLogMW(mux),
	}
	// Launching server in separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("HTTP Server stopped working")
		}
	}()

	// Graceful shutdown goroutine
	doneCh := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)
		for range sigCh {
			log.Println("Received an interrupt...")
			_ = server.Shutdown(ctx)
			doneCh <- struct{}{}
		}
	}()

	<-doneCh
}
