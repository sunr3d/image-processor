package entrypoint

import (
	"context"
	"strings"

	"github.com/sunr3d/image-processor/internal/config"
	httphandlers "github.com/sunr3d/image-processor/internal/handlers"
	"github.com/sunr3d/image-processor/internal/infra/broker/kafka"
	"github.com/sunr3d/image-processor/internal/infra/storage/filestorage"
	"github.com/sunr3d/image-processor/internal/server"
	"github.com/sunr3d/image-processor/internal/services/imagesvc"
)

func RunApp(ctx context.Context, cfg *config.Config) error {
	// Инфраслой (Infrastructure layer)
	imageStor := filestorage.NewFileStorage(cfg.StoragePath)
	metadataStor := filestorage.NewMetadataStorage(cfg.MetadataPath)

	kafkaBrokers := strings.Split(cfg.KafkaBrokers, ",")
	publisher := kafka.NewPublisher(kafkaBrokers, cfg.KafkaTopic)
	defer publisher.Close()

	// Сервисный слой (Application / Use Cases layer)
	imageSvc := imagesvc.New(imageStor, metadataStor, publisher)

	// Слой представления (Presentation layer)
	h := httphandlers.New(imageSvc)
	engine := h.RegisterHandlers()

	// Сервер
	srv := server.New(":"+cfg.HTTPPort, engine)
	
	return srv.Run(ctx)
}
