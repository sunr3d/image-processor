package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/image-processor/internal/interfaces/infra"
	"github.com/sunr3d/image-processor/internal/interfaces/services"
	"github.com/sunr3d/image-processor/models"
)

type worker struct {
	processor   services.ImageProcessor
	imgStorage  infra.ImageStorage
	metaStorage infra.MetadataStorage
	subscriber  infra.Subscriber
}

// New - конструктор Worker.
func New(proc services.ImageProcessor, imgStor infra.ImageStorage, metaStor infra.MetadataStorage, sub infra.Subscriber) *worker {
	return &worker{
		processor:   proc,
		imgStorage:  imgStor,
		metaStorage: metaStor,
		subscriber:  sub,
	}
}

func (w *worker) Start(ctx context.Context) error {
	zlog.Logger.Info().Msg("Worker запущен, ожидание задач из брокера...")

	return w.subscriber.Subscribe(ctx, w.processTask)
}

func (w *worker) processTask(ctx context.Context, task *models.ProcessingTask) error {
	zlog.Logger.Info().Msgf("Обработка задачи: %s", task.ImageID)

	meta, err := w.setMetaToProcessing(ctx, task.ImageID)
	if err != nil {
		return fmt.Errorf("setMetaToProcessing: %w", err)
	}

	result, err := w.processImg(task.OriginalPath)
	if err != nil {
		w.handleProcessingErr(ctx, meta, err)
		return fmt.Errorf("processImage: %w", err)
	}

	paths, err := w.saveImages(ctx, task.ImageID, result)
	if err != nil {
		w.handleProcessingErr(ctx, meta, err)
		return fmt.Errorf("saveImages: %w", err)
	}

	if err := w.setMetaToCompleted(ctx, meta, paths); err != nil {
		return fmt.Errorf("setMetaToCompleted: %w", err)
	}

	zlog.Logger.Info().Msgf("Задача %s успешно обработана", task.ImageID)

	return nil
}

// helpers
func (w *worker) setMetaToProcessing(ctx context.Context, imageID string) (*models.ImageMetadata, error) {
	meta, err := w.metaStorage.Get(ctx, imageID)
	if err != nil {
		return nil, fmt.Errorf("metaStorage.Get: %w", err)
	}

	meta.Status = models.StatusProcessing
	meta.UpdatedAt = time.Now()

	if err := w.metaStorage.Update(ctx, meta); err != nil {
		return nil, fmt.Errorf("metaStorage.Update: %w", err)
	}

	return meta, nil
}

func (w *worker) processImg(imagePath string) (*models.ProcessedImages, error) {
	result, err := w.processor.Process(imagePath)
	if err != nil {
		return nil, fmt.Errorf("processor.Process: %w", err)
	}

	return result, nil
}

func (w *worker) saveImages(ctx context.Context, imageID string, result *models.ProcessedImages) (map[string]string, error) {
	paths := make(map[string]string)

	// 1. Resized
	resizedPath, err := w.imgStorage.SaveProcessed(ctx, imageID, "resized", result.Resized)
	if err != nil {
		return nil, fmt.Errorf("imgStorage.SaveProcessed Resized: %w", err)
	}
	paths["resized"] = resizedPath

	// 2. Thumbnail
	thumbnailPath, err := w.imgStorage.SaveProcessed(ctx, imageID, "thumbnail", result.Thumbnail)
	if err != nil {
		return nil, fmt.Errorf("imgStorage.SaveProcessed Thumbnail: %w", err)
	}
	paths["thumbnail"] = thumbnailPath

	// 3. Watermarked
	watermarkedPath, err := w.imgStorage.SaveProcessed(ctx, imageID, "watermarked", result.Watermarked)
	if err != nil {
		return nil, fmt.Errorf("imgStorage.SaveProcessed Watermarked: %w", err)
	}
	paths["watermarked"] = watermarkedPath

	return paths, nil
}

func (w *worker) setMetaToCompleted(ctx context.Context, meta *models.ImageMetadata, paths map[string]string) error {
	meta.Status = models.StatusCompleted
	meta.ResizedPath = paths["resized"]
	meta.ThumbnailPath = paths["thumbnail"]
	meta.WatermarkedPath = paths["watermarked"]
	meta.UpdatedAt = time.Now()

	if err := w.metaStorage.Update(ctx, meta); err != nil {
		return fmt.Errorf("metaStorage.Update: %w", err)
	}

	return nil
}

func (w *worker) handleProcessingErr(ctx context.Context, meta *models.ImageMetadata, procErr error) error {
	meta.Status = models.StatusFailed
	meta.ErrorMessage = procErr.Error()
	meta.UpdatedAt = time.Now()

	if err := w.metaStorage.Update(ctx, meta); err != nil {
		return fmt.Errorf("metaStorage.Update: %w", err)
	}

	return nil
}
