package entrypoint

import (
	"context"
	"strings"

	"github.com/sunr3d/image-processor/internal/config"
	httphandlers "github.com/sunr3d/image-processor/internal/handlers"
	"github.com/sunr3d/image-processor/internal/infra/broker/kafka"
	"github.com/sunr3d/image-processor/internal/infra/storage/filestorage"
	"github.com/sunr3d/image-processor/internal/services/imagesvc"
)

func RunApp(ctx context.Context, cfg *config.Config) error {
	// Инфраслой
	imageStor := filestorage.NewFileStorage(cfg.StoragePath)
	metadataStor := filestorage.NewMetadataStorage(cfg.MetadataPath)

	kafkaBrokers := strings.Split(cfg.KafkaBrokers, ",")
	publisher := kafka.NewPublisher(kafkaBrokers, cfg.KafkaTopic)
	defer publisher.Close()

	// Сервисный слой
	imageSvc := imagesvc.New(imageStor, metadataStor, publisher)

	// REST API (HTTP) + Middleware
	h := httphandlers.New(imageSvc)
	engine := h.RegisterHandlers()

	// Server
	return engine.Run(":" + cfg.HTTPPort)
}
