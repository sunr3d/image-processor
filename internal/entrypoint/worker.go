package entrypoint

import (
	"context"
	"strings"

	"github.com/sunr3d/image-processor/internal/config"
	"github.com/sunr3d/image-processor/internal/infra/broker/kafka"
	"github.com/sunr3d/image-processor/internal/infra/storage/filestorage"
	"github.com/sunr3d/image-processor/internal/services/processor"
	"github.com/sunr3d/image-processor/internal/services/worker"
	"github.com/wb-go/wbf/zlog"
)

func RunWorker(ctx context.Context, cfg *config.Config) error {
	// Инфраслой
	imgStor := filestorage.NewFileStorage(cfg.StoragePath)
	metaStor := filestorage.NewMetadataStorage(cfg.MetadataPath)

	kafkaBrokers := strings.Split(cfg.KafkaBrokers, ",")
	subscriber := kafka.NewSubscriber(kafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroup)
	defer subscriber.Close()

	// Сервисный слой
	proc := processor.New(cfg.ThumbnailSize, cfg.ResizeWidth, cfg.WatermarkText)

	workerSvc := worker.New(proc, imgStor, metaStor, subscriber)

	zlog.Logger.Info().Msg("Worker запущен и ожидает задачи из Kafka")

	return workerSvc.Start(ctx)
}
