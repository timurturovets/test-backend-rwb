package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/timurturovets/test-backend-rwb/config"
	"github.com/timurturovets/test-backend-rwb/internal/api"
	"github.com/timurturovets/test-backend-rwb/internal/consumer"
	"github.com/timurturovets/test-backend-rwb/internal/engine"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	eng := engine.NewEngine()
	go eng.Start(ctx)

	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		logger.Error("failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer nc.Close()

	js, err := setupJetStream(nc)
	if err != nil {
		logger.Error("failed to setup JetStream", "error", err)
		os.Exit(1)
	}
	if err := consumer.EnsureStream(ctx, js, cfg.NATSStream, cfg.NATSSubject); err != nil {
		logger.Error("failed to ensure stream", "error", err)
		os.Exit(1)
	}

	cons, err := consumer.New(nc, cfg.NATSStream, cfg.NATSSubject,
		func(event consumer.SearchEvent) {
			eng.Add(event.Query)
		},
		logger,
	)
	if err != nil {
		logger.Error("failed to create consumer", "error", err)
		os.Exit(1)
	}
	go func() {
		if err := cons.Start(ctx); err != nil {
			logger.Error("consumer error", "error", err)
		}
	}()

	handler := api.NewHandler(eng)
	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		logger.Info("server starting", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")
	srv.Shutdown(context.Background())
}

func setupJetStream(nc *nats.Conn) (jetstream.JetStream, error) {
	return jetstream.New(nc)
}
