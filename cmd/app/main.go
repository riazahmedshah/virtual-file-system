package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/riazahmedshah/vfs/internal/config"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/repository"
	"github.com/riazahmedshah/vfs/internal/router"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/service"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config" + err.Error())
	}

	srv, err := server.New(cfg)
	if err != nil {
		slog.Error("failed to creat server", "err", err.Error())
	}

	repos := repository.NewRepositories(srv)
	services, serviceErr := service.NewServices(srv, repos)
	if serviceErr != nil {
		slog.Error("could not create services", "err", err)
		os.Exit(1)
	}
	handlers := handler.NewHandlers(srv, services)

	r := router.NewRouter(handlers)

	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		if err := srv.Start(); err != nil && errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced shut down", "err", err)
		os.Exit(1)
	}
	stop()
	cancel()

	slog.Info("server exited gracefully")
}
