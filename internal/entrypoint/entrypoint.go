package entrypoint

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunr3d/image-processor/internal/config"
	httphandlers "github.com/sunr3d/image-processor/internal/handlers"
)

func Run(cfg *config.Config) error {
	appCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Инфраслой
	// TODO: Implement infrastructure

	// Сервисный слой
	// TODO: Implement service

	// REST API (HTTP) + Middleware
	h := httphandlers.New(svc)
	engine := h.RegisterHandlers()

	// Server
	return engine.Run(":" + cfg.HTTPPort)
}
