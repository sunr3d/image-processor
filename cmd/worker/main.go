package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/config"
	"github.com/sunr3d/image-processor/internal/entrypoint"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("Запуск Worker сервиса...")

	cfg, err := config.GetConfig("config.yml")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("config.GetConfig")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := entrypoint.RunWorker(ctx, cfg); err != nil {
		zlog.Logger.Fatal().Err(err).Msg("entrypoint.RunWorker")
	}
}
